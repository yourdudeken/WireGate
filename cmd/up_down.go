package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/config"
)

var upCmd = &cobra.Command{
	Use:   "up [peerName]",
	Short: "Start a local peer (home server) tunnel",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		peerName := args[0]
		
		_, err := config.LoadConfig("config.yaml")
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		peerPath := filepath.Join("deploy/peers", peerName)
		if _, err := os.Stat(peerPath); os.IsNotExist(err) {
			fmt.Printf("Configurations for peer '%s' not found. Run 'generate' first.\n", peerName)
			return
		}

		fmt.Printf("Starting tunnel for peer '%s'...\n", peerName)
		
		dockerCmd := exec.Command("docker", "compose", "up", "-d")
		dockerCmd.Dir = peerPath
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr
		
		if err := dockerCmd.Run(); err != nil {
			fmt.Printf("Error starting peer: %v\n", err)
			return
		}

		fmt.Printf("Peer '%s' is now UP.\n", peerName)
	},
}

var downCmd = &cobra.Command{
	Use:   "down [peerName]",
	Short: "Stop a local peer (home server) tunnel",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		peerName := args[0]
		
		peerPath := filepath.Join("deploy/peers", peerName)
		if _, err := os.Stat(peerPath); os.IsNotExist(err) {
			fmt.Printf("Configurations for peer '%s' not found.\n", peerName)
			return
		}

		fmt.Printf("Stopping tunnel for peer '%s'...\n", peerName)
		
		dockerCmd := exec.Command("docker", "compose", "down")
		dockerCmd.Dir = peerPath
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr
		
		if err := dockerCmd.Run(); err != nil {
			fmt.Printf("Error stopping peer: %v\n", err)
			return
		}

		fmt.Printf("Peer '%s' is now DOWN.\n", peerName)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
}
