package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/config"
)

var sshCmd = &cobra.Command{
	Use:   "ssh [peer-name]",
	Short: "SSH into a peer (home server) via the WireGuard tunnel",
	Long: `Establishes an SSH connection to a peer using its internal WireGuard IP.
This command is intended to be run from the VPS or any machine connected to the same WireGuard network.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		var targetPeer *config.PeerConfig
		for i, p := range cfg.Peers {
			if p.Name == name {
				targetPeer = &cfg.Peers[i]
				break
			}
		}

		if targetPeer == nil {
			fmt.Printf("Peer '%s' not found.\n", name)
			return
		}

		user := targetPeer.SSHUser
		if user == "" {
			user = "root"
		}

		fmt.Printf("Connecting to %s (%s) as %s...\n", name, targetPeer.WGIp, user)
		
		sshArgs := []string{fmt.Sprintf("%s@%s", user, targetPeer.WGIp)}
		
		// Disable strict host key checking for internal WireGuard IPs to simplify access
		sshArgs = append([]string{"-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null"}, sshArgs...)

		sshCmd := exec.Command("ssh", sshArgs...)
		sshCmd.Stdout = os.Stdout
		sshCmd.Stderr = os.Stderr
		sshCmd.Stdin = os.Stdin

		if err := sshCmd.Run(); err != nil {
			// Do not print error if it was a clean exit
			if exitError, ok := err.(*exec.ExitError); ok {
				if exitError.ExitCode() != 0 {
					fmt.Printf("SSH connection failed: %v\n", err)
				}
			} else {
				fmt.Printf("Error: %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
