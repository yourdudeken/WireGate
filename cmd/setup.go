package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure the local environment for the gateway",
	Long: `Performs necessary local environment setup including:
- Configuring the local firewall (UFW) to permit tunnel traffic
- Verifying Docker installation
- Providing system-wide installation instructions`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("WireGate Local Setup")
		fmt.Println("------------------------")

		if runtime.GOOS != "linux" {
			fmt.Println("Notice: Automate firewall configuration is only supported on Linux (UFW).")
			fmt.Println("Please manually ensure your firewall allows traffic on the 'wg0' interface.")
		} else {
			fmt.Println("Configuring local firewall (UFW)...")
			fmt.Println("Command: sudo ufw allow in on wg0")
			
			fwCmd := exec.Command("sudo", "ufw", "allow", "in", "on", "wg0")
			fwCmd.Stdout = os.Stdout
			fwCmd.Stderr = os.Stderr
			if err := fwCmd.Run(); err != nil {
				fmt.Printf("Warning: Could not configure firewall automatically: %v\n", err)
				fmt.Println("Please run 'sudo ufw allow in on wg0' manually.")
			} else {
				fmt.Println("Firewall configured successfully.")
			}
		}

		fmt.Println("\nChecking Docker Engine...")
		dockerCmd := exec.Command("docker", "info")
		if err := dockerCmd.Run(); err != nil {
			fmt.Println("Error: Docker is not running or not installed.")
			fmt.Println("Please install Docker: https://docs.docker.com/get-docker/")
		} else {
			fmt.Println("Docker Engine: OK")
		}

		fmt.Println("\nSystem Installation:")
		fmt.Println("To run 'wiregate' from anywhere, run:")
		cwd, _ := os.Getwd()
		fmt.Printf("sudo ln -sf %s/wiregate /usr/local/bin/wiregate\n", cwd)
		
		fmt.Println("\nSetup complete. You can now use './wiregate up [peer]' to start your tunnels.")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
