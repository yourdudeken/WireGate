package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/config"
)

var followFlag bool

var logsCmd = &cobra.Command{
	Use:   "logs [vps|home] [service]",
	Short: "View logs from VPS or Home server",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		service := ""
		if len(args) > 1 {
			service = args[1]
		}

		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if target == "vps" {
			viewVPSLogs(cfg, service)
		} else {
			viewHomeLogs(cfg, target)
		}
	},
}

func viewVPSLogs(cfg *config.Config, service string) {
	dest := fmt.Sprintf("%s@%s", cfg.VPS.SSHUser, cfg.VPS.IP)
	
	remoteCmd := "cd ~/wiregate && docker compose logs"
	if followFlag {
		remoteCmd += " -f"
	}
	if service != "" {
		remoteCmd += " " + service
	}

	sshArgs := []string{dest, remoteCmd}
	if cfg.VPS.SSHKey != "" {
		sshArgs = append([]string{"-i", cfg.VPS.SSHKey}, sshArgs...)
	}

	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error viewing logs: %v\n", err)
	}
}

func viewHomeLogs(cfg *config.Config, target string) {
	fmt.Printf("Viewing logs for peer '%s'...\n", target)
	
	args := []string{"compose", "logs"}
	if followFlag {
		args = append(args, "-f")
	}

	cmd := exec.Command("docker", args...)
	cmd.Dir = fmt.Sprintf("deploy/peers/%s", target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error viewing logs: %v\n", err)
	}
}

func init() {
	logsCmd.Flags().BoolVarP(&followFlag, "follow", "f", false, "Follow log output")
	rootCmd.AddCommand(logsCmd)
}
