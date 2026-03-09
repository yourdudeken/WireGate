package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/config"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current configuration status",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig("config.yaml")
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		fmt.Println("WireGate Status")
		fmt.Println("======================")
		fmt.Printf("Project:    %s\n", cfg.Project)
		fmt.Printf("VPS IP:     %s (User: %s)\n", cfg.VPS.IP, cfg.VPS.SSHUser)
		fmt.Printf("Audit:      ")
		if err := cfg.Validate(); err != nil {
			fmt.Printf("Incomplete: %v\n", err)
		} else {
			fmt.Println("Production Ready")
		}
		
		fmt.Println("\nNetworking")
		fmt.Printf("  Tunnel IP (VPS):  %s\n", cfg.VPS.WGIp)
		fmt.Println("\nPeers:")
		if len(cfg.Peers) == 0 {
			fmt.Println("  No peers configured.")
		} else {
			for _, p := range cfg.Peers {
				fmt.Printf("  - %s: %s\n", p.Name, p.WGIp)
			}
		}
		
		fmt.Println("\nServices")
		if len(cfg.Services) == 0 {
			fmt.Println("  No services configured.")
		} else {
			for i, s := range cfg.Services {
				fmt.Printf("  %d. %s -> %s (Peer: %s)\n", i+1, s.Domain, "localhost", s.PeerName)
			}
		}
		fmt.Println("======================")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
