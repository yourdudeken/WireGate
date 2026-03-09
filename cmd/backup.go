package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/backup"
	"github.com/yourdudeken/wiregate/internal/config"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a full backup of configuration and certificates",
	Long: fmt.Sprintf("Creates a timestamped ZIP archive containing:\n- %s (Local machine)\n- LetsEncrypt certificates (Fetched from VPS)\nIf configured, the backup will be saved locally or uploaded to S3.", ConfigFile),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		fmt.Println("WireGate Backup System")
		fmt.Println("-------------------------")

		archivePath, err := backup.Run(cfg)
		if err != nil {
			fmt.Printf("Error during backup: %v\n", err)
			return
		}

		fmt.Printf("\nBackup created successfully: %s\n", archivePath)
		
		// Clean up temporary file if it was only intended for S3
		if cfg.Backup.LocalPath == "" && cfg.Backup.S3.Enabled {
			defer os.Remove(archivePath)
		}

		if cfg.Backup.S3.Enabled {
			fmt.Println("Uploading to S3...")
			// TODO: Implement S3 Upload Logic
			fmt.Println("S3 upload functionality coming soon in next patch (Use local backups for now)")
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
