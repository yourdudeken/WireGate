package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveConfig(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func NewDefaultConfig() *Config {
	return &Config{
		Project: "home-gateway",
		VPS: VPSConfig{
			WGPort: 51820,
			WGIp:   "10.0.0.1",
		},
		Peers: []PeerConfig{
			{
				Name:      "home",
				WGIp:      "10.0.0.2",
				SSHUser:   "root",
				Keepalive: 25,
			},
		},
		Proxy: ProxyConfig{
			Type: "traefik",
		},
		Monitor: MonitorConfig{
			Interval: 5,
		},
	}
}
