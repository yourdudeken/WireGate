package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/yourdudeken/wiregate/internal/config"
	"github.com/yourdudeken/wiregate/internal/ssh"
)

type Hub struct {
	config *config.Config
	client *ssh.Client
}

func NewHub(cfg *config.Config) *Hub {
	client := ssh.NewClient(cfg.VPS.SSHUser, cfg.VPS.IP, cfg.VPS.SSHKey)
	return &Hub{
		config: cfg,
		client: client,
	}
}

func (h *Hub) Start() {
	interval := h.config.Monitor.Interval
	if interval <= 0 {
		interval = 5 // Default to 5 minutes
	}

	log.Printf("Starting Monitor (Interval: %d minutes)\n", interval)
	
	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	
	// Run initially
	h.Check()

	for range ticker.C {
		h.Check()
	}
}

func (h *Hub) Check() {
	log.Println("Running Proactive Health Check...")
	
	// 1. VPS SSH Check
	if err := h.client.Run("true"); err != nil {
		h.alert("CRITICAL", fmt.Sprintf("VPS at %s is UNREACHABLE via SSH", h.config.VPS.IP))
		return
	}

	// 2. Peer Checks
	for _, peer := range h.config.Peers {
		if err := h.client.Run(fmt.Sprintf("docker exec wireguard ping -c 2 -W 1 %s > /dev/null 2>&1", peer.WGIp)); err != nil {
			h.alert("WARNING", fmt.Sprintf("Peer '%s' (%s) is OFFLINE", peer.Name, peer.WGIp))
		}
	}

	// 3. Service Checks
	for _, svc := range h.config.Services {
		var peer config.PeerConfig
		found := false
		for _, p := range h.config.Peers {
			if p.Name == svc.PeerName {
				peer = p
				found = true
				break
			}
		}
		
		if found {
			if err := h.client.Run(fmt.Sprintf("docker exec wireguard nc -z -w 3 %s %d > /dev/null 2>&1", peer.WGIp, svc.Port)); err != nil {
				h.alert("CRITICAL", fmt.Sprintf("Service '%s' (%s:%d) is NOT RESPONDING", svc.Domain, peer.WGIp, svc.Port))
			}
		}
	}
}

func (h *Hub) alert(level, message string) {
	log.Printf("[%s] %s\n", level, message)
	
	alert := Alert{
		Level:   level,
		Message: message,
		Time:    time.Now(),
	}

	if h.config.Monitor.Discord.Enabled {
		if err := SendDiscord(h.config.Monitor.Discord.URL, alert); err != nil {
			log.Printf("Failed to send Discord alert: %v\n", err)
		}
	}

	if h.config.Monitor.Telegram.Enabled {
		if err := SendTelegram(h.config.Monitor.Telegram.Token, h.config.Monitor.Telegram.ChatID, alert); err != nil {
			log.Printf("Failed to send Telegram alert: %v\n", err)
		}
	}
}
