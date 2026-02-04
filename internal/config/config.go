package config

type Config struct {
	Project  string        `yaml:"project"`
	VPS      VPSConfig     `yaml:"vps"`
	Peers    []PeerConfig  `yaml:"peers"`
	Proxy    ProxyConfig   `yaml:"proxy"`
	Services []Service     `yaml:"services"`
	Monitor  MonitorConfig `yaml:"monitor"`
	Backup   BackupConfig  `yaml:"backup"`
}

type MonitorConfig struct {
	Interval int           `yaml:"interval"` // in minutes
	Discord  WebHookConfig `yaml:"discord"`
	Telegram WebHookConfig `yaml:"telegram"`
}

type BackupConfig struct {
	LocalPath string   `yaml:"local_path"`
	S3        S3Config `yaml:"s3"`
}

type S3Config struct {
	Enabled   bool   `yaml:"enabled"`
	Endpoint  string `yaml:"endpoint"`
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

type WebHookConfig struct {
	Enabled bool   `yaml:"enabled"`
	URL     string `yaml:"url"`
	Token   string `yaml:"token"`   // for telegram
	ChatID  string `yaml:"chat_id"` // for telegram
}

type VPSConfig struct {
	IP         string `yaml:"ip"`
	SSHUser    string `yaml:"ssh_user"`
	SSHKey     string `yaml:"ssh_key"`
	WGIp       string `yaml:"wg_ip"`
	WGPort     int    `yaml:"wg_port"`
	PrivateKey string `yaml:"private_key"`
	PublicKey  string `yaml:"public_key"`
}

type PeerConfig struct {
	Name       string `yaml:"name"`
	WGIp       string `yaml:"wg_ip"`
	SSHUser    string `yaml:"ssh_user"`
	Keepalive  int    `yaml:"keepalive"`
	PrivateKey string `yaml:"private_key"`
	PublicKey  string `yaml:"public_key"`
}

type ProxyConfig struct {
	Type  string `yaml:"type"` // e.g. "traefik"
	Email string `yaml:"email"`
}

type Service struct {
	Name     string `yaml:"name"`
	Domain   string `yaml:"domain"`
	Port     int    `yaml:"port"`
	PeerName string `yaml:"peer_name"`
}
