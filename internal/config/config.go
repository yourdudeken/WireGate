package config

type Config struct {
	Project  string        `yaml:"project" json:"project"`
	VPS      VPSConfig     `yaml:"vps" json:"vps"`
	Peers    []PeerConfig  `yaml:"peers" json:"peers"`
	Proxy    ProxyConfig   `yaml:"proxy" json:"proxy"`
	Services []Service     `yaml:"services" json:"services"`
	Monitor  MonitorConfig `yaml:"monitor" json:"monitor"`
	Backup   BackupConfig  `yaml:"backup" json:"backup"`
}

type MonitorConfig struct {
	Interval int           `yaml:"interval" json:"interval"` // in minutes
	Discord  WebHookConfig `yaml:"discord" json:"discord"`
	Telegram WebHookConfig `yaml:"telegram" json:"telegram"`
}

type BackupConfig struct {
	LocalPath string   `yaml:"local_path" json:"local_path"`
	S3        S3Config `yaml:"s3" json:"s3"`
}

type S3Config struct {
	Enabled   bool   `yaml:"enabled" json:"enabled"`
	Endpoint  string `yaml:"endpoint" json:"endpoint"`
	Region    string `yaml:"region" json:"region"`
	Bucket    string `yaml:"bucket" json:"bucket"`
	AccessKey string `yaml:"access_key" json:"access_key"`
	SecretKey string `yaml:"secret_key" json:"secret_key"`
}

type WebHookConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	URL     string `yaml:"url" json:"url"`
	Token   string `yaml:"token" json:"token"`   // for telegram
	ChatID  string `yaml:"chat_id" json:"chat_id"` // for telegram
}

type VPSConfig struct {
	IP         string `yaml:"ip" json:"ip"`
	SSHUser    string `yaml:"ssh_user" json:"ssh_user"`
	SSHKey     string `yaml:"ssh_key" json:"ssh_key"`
	WGIp       string `yaml:"wg_ip" json:"wg_ip"`
	WGPort     int    `yaml:"wg_port" json:"wg_port"`
	PrivateKey string `yaml:"private_key" json:"private_key"`
	PublicKey  string `yaml:"public_key" json:"public_key"`
}

type PeerConfig struct {
	Name       string `yaml:"name" json:"name"`
	WGIp       string `yaml:"wg_ip" json:"wg_ip"`
	SSHUser    string `yaml:"ssh_user" json:"ssh_user"`
	Keepalive  int    `yaml:"keepalive" json:"keepalive"`
	PrivateKey string `yaml:"private_key" json:"private_key"`
	PublicKey  string `yaml:"public_key" json:"public_key"`
}

type ProxyConfig struct {
	Type  string `yaml:"type" json:"type"` // e.g. "traefik"
	Email string `yaml:"email" json:"email"`
}

type Service struct {
	Name     string `yaml:"name" json:"name"`
	Domain   string `yaml:"domain" json:"domain"`
	Port     int    `yaml:"port" json:"port"`
	PeerName string `yaml:"peer_name" json:"peer_name"`
}
