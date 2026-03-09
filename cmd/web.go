package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/web"
)

var (
	webPort     int
	webPassword string
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web UI dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		// Check environment variable first
		if envPass := os.Getenv("WG_ADMIN_PASS"); envPass != "" {
			webPassword = envPass
		}

		server := web.NewServer(ConfigFile, webPassword)
		
		fmt.Printf("Starting WireGate Web UI...\n")
		fmt.Printf("Dashboard: http://localhost:%d\n", webPort)
		if webPassword != "" {
			fmt.Printf("Authentication: ENABLED (User: admin)\n")
		} else {
			fmt.Printf("Authentication: DISABLED (Warning: Dashboard is public!)\n")
		}
		fmt.Printf("Press Ctrl+C to stop\n\n")
		
		if err := server.Start(webPort); err != nil {
			fmt.Printf("Error starting web server: %v\n", err)
		}
	},
}

func init() {
	webCmd.Flags().IntVarP(&webPort, "port", "p", 8080, "Port to run the web UI on")
	webCmd.Flags().StringVar(&webPassword, "password", "", "Password for dashboard authentication (can also use WG_ADMIN_PASS env var)")
	rootCmd.AddCommand(webCmd)
}
