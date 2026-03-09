package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/config"
	"github.com/yourdudeken/wiregate/internal/service"
)

var (
	targetPeer string
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage home services",
}

var addSvcCmd = &cobra.Command{
	Use:   "add [domain] [port]",
	Short: "Add a new service",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		port, _ := strconv.Atoi(args[1])

		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		// If domain doesn't contain a dot, assume it's a prefix and use sslip.io
		if !strings.Contains(domain, ".") {
			domain = fmt.Sprintf("%s.%s.sslip.io", domain, cfg.VPS.IP)
		}

		if err := service.Validate(domain, port); err != nil {
			fmt.Printf("Validation error: %v\n", err)
			return
		}

		// Use domain as name if not specified
		name := domain

		if err := service.Add(cfg, name, domain, port, targetPeer); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := config.SaveConfig(ConfigFile, cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Service %s added successfully to peer %s.\n", domain, targetPeer)
	},
}

var updateSvcCmd = &cobra.Command{
	Use:   "update [domain] [new-port]",
	Short: "Update an existing service's port",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		newPort, _ := strconv.Atoi(args[1])

		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if err := service.Edit(cfg, domain, newPort); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := config.SaveConfig(ConfigFile, cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Service %s updated to port %d.\n", domain, newPort)
	},
}

var removeSvcCmd = &cobra.Command{
	Use:   "remove [domain]",
	Short: "Remove an existing service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]

		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if err := service.Remove(cfg, domain); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := config.SaveConfig(ConfigFile, cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Service %s removed successfully.\n", domain)
	},
}

var listSvcCmd = &cobra.Command{
	Use:   "list",
	Short: "List all services",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if len(cfg.Services) == 0 {
			fmt.Println("No services configured.")
			return
		}

		fmt.Println("Current Services:")
		for _, s := range cfg.Services {
			fmt.Printf("  - %s -> %s:%d (Peer: %s)\n", s.Domain, "localhost", s.Port, s.PeerName)
		}
	},
}

var addTemplateCmd = &cobra.Command{
	Use:   "add-template [template-name] [domain-prefix]",
	Short: "Add a new service using a pre-configured template",
	Long: `Available templates: plex, ha, jellyfin, pihole, grafana, prometheus, overseerr, tautulli.
Example: wiregate service add-template plex myplex`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tmplName := strings.ToLower(args[0])
		prefix := args[1]

		tmpl, ok := service.Templates[tmplName]
		if !ok {
			fmt.Printf("Error: Template '%s' not found.\n", tmplName)
			return
		}

		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		// Auto-generate domain if it's just a prefix
		domain := prefix
		if !strings.Contains(domain, ".") {
			domain = fmt.Sprintf("%s.%s.sslip.io", prefix, cfg.VPS.IP)
		}

		if err := service.Add(cfg, domain, domain, tmpl.Port, targetPeer); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := config.SaveConfig(ConfigFile, cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Template '%s' applied: Service %s added to peer %s on port %d.\n", tmplName, domain, targetPeer, tmpl.Port)
	},
}

func init() {
	addSvcCmd.Flags().StringVar(&targetPeer, "peer", "home", "Target peer name for the service")
	addTemplateCmd.Flags().StringVar(&targetPeer, "peer", "home", "Target peer name for the service")
	serviceCmd.AddCommand(addSvcCmd)
	serviceCmd.AddCommand(updateSvcCmd)
	serviceCmd.AddCommand(removeSvcCmd)
	serviceCmd.AddCommand(listSvcCmd)
	serviceCmd.AddCommand(addTemplateCmd)
	rootCmd.AddCommand(serviceCmd)
}
