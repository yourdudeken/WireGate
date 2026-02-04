package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/yourdudeken/wg-gateway/internal/config"
	"github.com/yourdudeken/wg-gateway/internal/service"
	"github.com/yourdudeken/wg-gateway/internal/ssh"
	"github.com/yourdudeken/wg-gateway/internal/wg"
	remoteSSH "golang.org/x/crypto/ssh"
	"time"
)

//go:embed dist
var frontendContent embed.FS

type Server struct {
	configPath string
	password   string
	cache      *StatusCache
}

type StatusCache struct {
	LastUpdate time.Time
	WGData     map[string]PeerLiveStatus
}

type PeerLiveStatus struct {
	Handshake  uint64 `json:"handshake"`
	TransferRx uint64 `json:"rx"`
	TransferTx uint64 `json:"tx"`
	Online     bool   `json:"online"`
}

func NewServer(configPath, password string) *Server {
	return &Server{
		configPath: configPath,
		password:   password,
		cache: &StatusCache{
			WGData: make(map[string]PeerLiveStatus),
		},
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
			fmt.Printf("Unauthorized access attempt from %s\n", r.RemoteAddr)
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

	s.updateCacheIfNeeded(cfg)
	var totalRx, totalTx uint64
	for _, p := range s.cache.WGData {
		totalRx += p.TransferRx
		totalTx += p.TransferTx
	}

	status := map[string]interface{}{
		"project":       cfg.Project,
		"vps_ip":        cfg.VPS.IP,
		"vps_user":      cfg.VPS.SSHUser,
		"ready":         cfg.Validate() == nil,
		"peer_count":    len(cfg.Peers),
		"service_count": len(cfg.Services),
		"total_rx":      totalRx,
		"total_tx":      totalTx,
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

	s.updateCacheIfNeeded(cfg)

	type PeerWithStatus struct {
		config.PeerConfig
		Live PeerLiveStatus `json:"live"`
	}

	peers := make([]PeerWithStatus, len(cfg.Peers))
	for i, p := range cfg.Peers {
		peers[i] = PeerWithStatus{
			PeerConfig: p,
			Live:       s.cache.WGData[p.PublicKey],
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}

func (s *Server) updateCacheIfNeeded(cfg *config.Config) {
	if time.Since(s.cache.LastUpdate) < 15*time.Second {
		return
	}

	if cfg.VPS.IP == "" {
		return
	}

	client := ssh.NewClient(cfg.VPS.SSHUser, cfg.VPS.IP, cfg.VPS.SSHKey)
	out, err := client.Output("sudo wg show all dump")
	if err != nil {
		return
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}
		// wg dump format: interface public_key preshared_key endpoint allowed_ips latest_handshake transfer_rx transfer_tx persistent_keepalive
		pubKey := fields[1]
		handshakeSec, _ := strconv.ParseInt(fields[5], 10, 64)
		rx, _ := strconv.ParseUint(fields[6], 10, 64)
		tx, _ := strconv.ParseUint(fields[7], 10, 64)

		s.cache.WGData[pubKey] = PeerLiveStatus{
			Handshake:  uint64(handshakeSec),
			TransferRx: rx,
			TransferTx: tx,
			Online:     handshakeSec > 0 && time.Now().Unix()-handshakeSec < 180, // online if handshake in last 3 mins
		}
	}
	s.cache.LastUpdate = time.Now()
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
	sshConfig := &remoteSSH.ClientConfig{
		User:            user,
		Auth:            []remoteSSH.AuthMethod{},
		HostKeyCallback: remoteSSH.InsecureIgnoreHostKey(),
	}

	// Try using SSH key
	if cfg.VPS.SSHKey != "" {
		key, err := os.ReadFile(cfg.VPS.SSHKey)
		if err == nil {
			signer, err := remoteSSH.ParsePrivateKey(key)
			if err == nil {
				sshConfig.Auth = append(sshConfig.Auth, remoteSSH.PublicKeys(signer))
			}
		}
	}

	// Connect to Peer via WG IP
	client, err := remoteSSH.Dial("tcp", fmt.Sprintf("%s:22", targetPeer.WGIp), sshConfig)
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
	modes := remoteSSH.TerminalModes{
		remoteSSH.ECHO:          1,
		remoteSSH.TTY_OP_ISPEED: 14400,
		remoteSSH.TTY_OP_OSPEED: 14400,
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
