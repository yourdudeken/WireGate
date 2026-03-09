package service

import (
	"errors"
	"fmt"
	"github.com/yourdudeken/wiregate/internal/config"
)

func Add(cfg *config.Config, name, domain string, port int, peerName string) error {
	for _, s := range cfg.Services {
		if s.Domain == domain {
			return fmt.Errorf("service with domain %s already exists", domain)
		}
	}

	// Verify peer exists
	peerExists := false
	for _, p := range cfg.Peers {
		if p.Name == peerName {
			peerExists = true
			break
		}
	}
	if !peerExists {
		return fmt.Errorf("peer %s not found", peerName)
	}

	cfg.Services = append(cfg.Services, config.Service{
		Name:     name,
		Domain:   domain,
		Port:     port,
		PeerName: peerName,
	})
	return nil
}

func Remove(cfg *config.Config, domain string) error {
	for i, s := range cfg.Services {
		if s.Domain == domain {
			cfg.Services = append(cfg.Services[:i], cfg.Services[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("service with domain %s not found", domain)
}

func Edit(cfg *config.Config, domain string, newPort int) error {
	for i, s := range cfg.Services {
		if s.Domain == domain {
			cfg.Services[i].Port = newPort
			return nil
		}
	}
	return fmt.Errorf("service with domain %s not found", domain)
}

func Validate(domain string, port int) error {
	if domain == "" {
		return errors.New("domain cannot be empty")
	}
	if port < 1 || port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	return nil
}
