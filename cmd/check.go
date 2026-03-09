package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/config"
	"github.com/yourdudeken/wiregate/internal/ssh"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the health and status of the gateway",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if err := cfg.Validate(); err != nil {
			fmt.Printf("Config is incomplete: %v\n", err)
			return
		}

		client := ssh.NewClient(cfg.VPS.SSHUser, cfg.VPS.IP, cfg.VPS.SSHKey)

		fmt.Println("Health Check: VPS Connection")
		fmt.Println("---------------------------")
		
		fmt.Print("SSH Connectivity: ")
		if err := client.Run("true"); err != nil {
			fmt.Println("FAILED")
			return
		}
		fmt.Println("OK")

		fmt.Print("Docker Engine:   ")
		if err := client.Run("docker info > /dev/null 2>&1"); err != nil {
			fmt.Println("FAILED")
		} else {
			fmt.Println("OK")
		}

		fmt.Print("WireGuard (wg0): ")
		if err := client.Run("docker exec wireguard ip addr show wg0 > /dev/null 2>&1"); err != nil {
			fmt.Println("DOWN")
		} else {
			fmt.Println("UP")
		}

		fmt.Println("\nPeer Connectivity:")
		for _, peer := range cfg.Peers {
			fmt.Printf("  - Peer %s (%s): ", peer.Name, peer.WGIp)
			if err := client.Run(fmt.Sprintf("docker exec wireguard ping -c 2 -W 1 %s > /dev/null 2>&1", peer.WGIp)); err != nil {
				fmt.Println("UNREACHABLE")
			} else {
				fmt.Println("REACHABLE")
			}
		}

		fmt.Println("\nService Connectivity:")
		for _, svc := range cfg.Services {
			var peer config.PeerConfig
			for _, p := range cfg.Peers {
				if p.Name == svc.PeerName {
					peer = p
					break
				}
			}
			fmt.Printf("  - %s (%s:%d): ", svc.Domain, peer.WGIp, svc.Port)
			if err := client.Run(fmt.Sprintf("docker exec wireguard nc -z -w 3 %s %d > /dev/null 2>&1", peer.WGIp, svc.Port)); err != nil {
				fmt.Println("TIMEOUT (Check local firewall)")
			} else {
				fmt.Println("OK")
			}
		}
		
		fmt.Println("---------------------------")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
