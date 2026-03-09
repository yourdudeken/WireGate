package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/config"
	"github.com/yourdudeken/wiregate/internal/wg"
)

var (
	peerWGIp    string
	peerSSHUser string
)

var peerCmd = &cobra.Command{
	Use:   "peer",
	Short: "Manage WireGuard peers (home servers)",
}

var addPeerCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new peer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		for _, p := range cfg.Peers {
			if p.Name == name {
				fmt.Printf("Peer %s already exists.\n", name)
				return
			}
		}

		// Validate WG IP if provided
		if peerWGIp != "" {
			if net.ParseIP(peerWGIp) == nil {
				fmt.Printf("Invalid WireGuard IP: %s\n", peerWGIp)
				return
			}
		} else {
			// Auto-assign next IP (e.g. 10.0.0.X)
			// Simple logic: find max last digit
			lastDigit := 2
			for _, p := range cfg.Peers {
				ip := net.ParseIP(p.WGIp)
				if ip != nil && len(ip) == 16 {
					if int(ip[15]) >= lastDigit {
						lastDigit = int(ip[15]) + 1
					}
				}
			}
			peerWGIp = fmt.Sprintf("10.0.0.%d", lastDigit)
		}

		keys, err := wg.GenerateKeyPair()
		if err != nil {
			fmt.Printf("Error generating keys: %v\n", err)
			return
		}

		cfg.Peers = append(cfg.Peers, config.PeerConfig{
			Name:       name,
			WGIp:       peerWGIp,
			SSHUser:    peerSSHUser,
			Keepalive:  25,
			PrivateKey: keys.Private,
			PublicKey:  keys.Public,
		})

		if err := config.SaveConfig(ConfigFile, cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Peer %s added successfully with IP %s.\n", name, peerWGIp)
		fmt.Println("Run 'generate' to create deployment files for this peer.")
	},
}

var listPeerCmd = &cobra.Command{
	Use:   "list",
	Short: "List all peers",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		fmt.Println("Current Peers:")
		for _, p := range cfg.Peers {
			fmt.Printf("  - %s: %s (Public Key: %s)\n", p.Name, p.WGIp, p.PublicKey)
		}
	},
}

func init() {
	addPeerCmd.Flags().StringVar(&peerWGIp, "ip", "", "WireGuard internal IP (optional)")
	addPeerCmd.Flags().StringVar(&peerSSHUser, "ssh-user", "root", "Default SSH user for this peer")
	peerCmd.AddCommand(addPeerCmd)
	peerCmd.AddCommand(listPeerCmd)
	rootCmd.AddCommand(peerCmd)
}
