package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/yourdudeken/wg-gateway/internal/config"
	"github.com/yourdudeken/wg-gateway/internal/service"
	"github.com/yourdudeken/wg-gateway/internal/wg"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

//go:embed dist
var frontendContent embed.FS

type Server struct {
	configPath string
	password   string
}

func NewServer(configPath, password string) *Server {
	return &Server{
		configPath: configPath,
		password:   password,
	}
}

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.password == "" {
			next(w, r)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != s.password {
			w.Header().Set("WWW-Authenticate", `Basic realm="W-G Gateway Dashboard"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func (s *Server) Start(port int) error {
	// API routes
	http.HandleFunc("/api/status", s.authMiddleware(s.handleStatus))
	http.HandleFunc("/api/peers", s.authMiddleware(s.handlePeers))
	http.HandleFunc("/api/peers/add", s.authMiddleware(s.handleAddPeer))
	http.HandleFunc("/api/services", s.authMiddleware(s.handleServices))
	http.HandleFunc("/api/services/add", s.authMiddleware(s.handleAddService))
	http.HandleFunc("/api/services/delete", s.authMiddleware(s.handleDeleteService))
	http.HandleFunc("/api/config", s.authMiddleware(s.handleConfig))
	http.HandleFunc("/api/config/update", s.authMiddleware(s.handleUpdateConfig))
	http.HandleFunc("/api/ssh", s.authMiddleware(s.handleSSH))

	// SPA serving
	distFS, _ := fs.Sub(frontendContent, "dist")
	fileServer := http.FileServer(http.FS(distFS))

	http.HandleFunc("/", s.authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Serve file if it exists, otherwise serve index.html for SPA routing
		f, err := distFS.Open(strings.TrimPrefix(r.URL.Path, "/"))
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}
		
		// Serve index.html
		indexFile, _ := distFS.Open("index.html")
		defer indexFile.Close()
		stat, _ := indexFile.Stat()
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
	}))

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig(s.configPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status := map[string]interface{}{
		"project":   cfg.Project,
		"vps_ip":    cfg.VPS.IP,
		"vps_user":  cfg.VPS.SSHUser,
		"ready":     cfg.Validate() == nil,
		"peer_count": len(cfg.Peers),
		"service_count": len(cfg.Services),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (s *Server) handlePeers(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig(s.configPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg.Peers)
}

func (s *Server) handleAddPeer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg, err := config.LoadConfig(s.configPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	keys, err := wg.GenerateKeyPair()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newPeer := config.PeerConfig{
		Name:       req.Name,
		WGIp:       req.IP,
		Keepalive:  25,
		PrivateKey: keys.Private,
		PublicKey:  keys.Public,
	}

	cfg.Peers = append(cfg.Peers, newPeer)

	if err := config.SaveConfig(s.configPath, cfg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newPeer)
}

func (s *Server) handleServices(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig(s.configPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg.Services)
}

func (s *Server) handleAddService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Domain   string `json:"domain"`
		Port     int    `json:"port"`
		PeerName string `json:"peer_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg, err := config.LoadConfig(s.configPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := service.Add(cfg, req.Domain, req.Domain, req.Port, req.PeerName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := config.SaveConfig(s.configPath, cfg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) handleDeleteService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "domain required", http.StatusBadRequest)
		return
	}

	cfg, err := config.LoadConfig(s.configPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := service.Remove(cfg, domain); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := config.SaveConfig(s.configPath, cfg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig(s.configPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func (s *Server) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg, err := config.LoadConfig(s.configPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch req.Key {
	case "vps.ip":
		cfg.VPS.IP = req.Value
	case "vps.user":
		cfg.VPS.SSHUser = req.Value
	case "proxy.email":
		cfg.Proxy.Email = req.Value
	case "project":
		cfg.Project = req.Value
	default:
		http.Error(w, "Unknown config key", http.StatusBadRequest)
		return
	}

	if err := config.SaveConfig(s.configPath, cfg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) handleSSH(w http.ResponseWriter, r *http.Request) {
	peerName := r.URL.Query().Get("peer")
	if peerName == "" {
		http.Error(w, "peer required", http.StatusBadRequest)
		return
	}

	cfg, err := config.LoadConfig(s.configPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var targetPeer *config.PeerConfig
	for i, p := range cfg.Peers {
		if p.Name == peerName {
			targetPeer = &cfg.Peers[i]
			break
		}
	}

	if targetPeer == nil {
		http.Error(w, "peer not found", http.StatusNotFound)
		return
	}

	user := targetPeer.SSHUser
	if user == "" {
		user = cfg.VPS.SSHUser
	}
	if user == "" {
		user = "root"
	}

	// Upgrade to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	// SSH Client Config
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Try using SSH key
	if cfg.VPS.SSHKey != "" {
		key, err := os.ReadFile(cfg.VPS.SSHKey)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
			}
		}
	}

	// Connect to Peer via WG IP
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", targetPeer.WGIp), sshConfig)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("\r\nConnection failed: "+err.Error()+"\r\n"))
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	// Request PTY
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm-256color", 40, 80, modes); err != nil {
		return
	}

	sshIn, _ := session.StdinPipe()
	sshOut, _ := session.StdoutPipe()

	go func() {
		io.Copy(sshIn, &wsWrapper{ws})
	}()

	go func() {
		io.Copy(&wsWrapper{ws}, sshOut)
	}()

	if err := session.Shell(); err != nil {
		return
	}

	session.Wait()
}

type wsWrapper struct {
	*websocket.Conn
}

func (w *wsWrapper) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(websocket.BinaryMessage, p)
	return len(p), err
}

func (w *wsWrapper) Read(p []byte) (n int, err error) {
	_, msg, err := w.ReadMessage()
	if err != nil {
		return 0, err
	}
	copy(p, msg)
	return len(msg), nil
}
