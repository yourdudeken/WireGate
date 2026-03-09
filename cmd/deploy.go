package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/config"
	"github.com/yourdudeken/wiregate/internal/provision"
	"github.com/yourdudeken/wiregate/internal/ssh"
)

var bootstrapFlag bool

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Automate VPS setup and deployment",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if cfg.VPS.IP == "" {
			fmt.Println("Error: VPS IP is not set. Use 'wiregate config vps.ip <ip>' first.")
			return
		}

		client := ssh.NewClient(cfg.VPS.SSHUser, cfg.VPS.IP, cfg.VPS.SSHKey)

		// 1. Optional Provisioning
		if bootstrapFlag {
			if err := provision.Bootstrap(client); err != nil {
				fmt.Printf("Provisioning failed: %v\n", err)
				return
			}
		}

		// 2. Generate local files
		fmt.Println("Generating deployment files...")
		generateAction(cfg)

		// 3. Deploy to VPS
		fmt.Printf("Deploying to VPS (%s)...\n", cfg.VPS.IP)
		
		// Create directory on VPS
		if err := client.Run("mkdir -p ~/wiregate/traefik_dynamic ~/wiregate/wireguard ~/wiregate/letsencrypt"); err != nil {
			fmt.Printf("Error creating directory on VPS: %v\n", err)
			return
		}

		// Upload files
		fmt.Println("Uploading configurations...")
		if err := client.Copy("deploy/vps/.", "~/wiregate"); err != nil {
			fmt.Printf("Error uploading files: %v\n", err)
			return
		}

		// Start services
		fmt.Println("Starting services on VPS...")
		if err := client.Run("cd ~/wiregate && docker compose up -d || docker-compose up -d"); err != nil {
			fmt.Printf("Error starting services: %v\n", err)
			return
		}

		fmt.Println("\nSuccess! Your VPS-to-Home Gateway is live.")
		fmt.Println("Now run 'docker compose up -d' in your 'deploy/home' directory on your home server.")
	},
}

func init() {
	deployCmd.Flags().BoolVarP(&bootstrapFlag, "bootstrap", "b", false, "Bootstrap the VPS (install Docker, WireGuard, etc.)")
	rootCmd.AddCommand(deployCmd)
}
