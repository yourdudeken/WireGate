package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/config"
	"github.com/yourdudeken/wiregate/internal/wg"
)

var (
	vpsIP      string
	sshUser    string
	sshKey     string
	proxyEmail string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(ConfigFile); err == nil {
			fmt.Printf("Config file '%s' already exists.\n", ConfigFile)
			return
		}

		cfg := config.NewDefaultConfig()

		if vpsIP != "" {
			cfg.VPS.IP = vpsIP
		}
		if sshUser != "" {
			cfg.VPS.SSHUser = sshUser
		}
		if sshKey != "" {
			cfg.VPS.SSHKey = sshKey
		}
		if proxyEmail != "" {
			cfg.Proxy.Email = proxyEmail
		}

		vpsKeys, err := wg.GenerateKeyPair()
		if err != nil {
			fmt.Printf("Error generating VPS keys: %v\n", err)
			return
		}
		cfg.VPS.PrivateKey = vpsKeys.Private
		cfg.VPS.PublicKey = vpsKeys.Public

		// Initialize first peer
		homeKeys, err := wg.GenerateKeyPair()
		if err != nil {
			fmt.Printf("Error generating peer keys: %v\n", err)
			return
		}
		cfg.Peers[0].PrivateKey = homeKeys.Private
		cfg.Peers[0].PublicKey = homeKeys.Public

		err = config.SaveConfig(ConfigFile, cfg)
		if err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Project initialized successfully at %s\n", ConfigFile)
		if cfg.VPS.IP == "" {
			fmt.Printf("Note: You still need to set your VPS IP using 'wiregate -c %s config vps.ip <ip>'\n", ConfigFile)
		}
	},
}

func init() {
	initCmd.Flags().StringVar(&vpsIP, "ip", "", "VPS Public IP address")
	initCmd.Flags().StringVar(&sshUser, "user", "root", "VPS SSH user")
	initCmd.Flags().StringVar(&sshKey, "key", "", "Path to SSH private key")
	initCmd.Flags().StringVar(&proxyEmail, "email", "", "Email for Let's Encrypt certificates")
	rootCmd.AddCommand(initCmd)
}
