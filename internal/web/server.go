package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/yourdudeken/wiregate/internal/config"
	"github.com/yourdudeken/wiregate/internal/service"
	"github.com/yourdudeken/wiregate/internal/wg"
)

//go:embed templates/* static/*
var content embed.FS

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
			w.Header().Set("WWW-Authenticate", `Basic realm="WireGate Dashboard"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func (s *Server) Start(port int) error {
	http.HandleFunc("/", s.authMiddleware(s.handleIndex))
	http.HandleFunc("/api/status", s.authMiddleware(s.handleStatus))
	http.HandleFunc("/api/peers", s.authMiddleware(s.handlePeers))
	http.HandleFunc("/api/peers/add", s.authMiddleware(s.handleAddPeer))
	http.HandleFunc("/api/services", s.authMiddleware(s.handleServices))
	http.HandleFunc("/api/services/add", s.authMiddleware(s.handleAddService))
	http.HandleFunc("/api/services/delete", s.authMiddleware(s.handleDeleteService))
	http.HandleFunc("/api/config", s.authMiddleware(s.handleConfig))
	http.HandleFunc("/api/config/update", s.authMiddleware(s.handleUpdateConfig))

	http.Handle("/static/", s.authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.FS(content)).ServeHTTP(w, r)
	}))

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(content, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
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
