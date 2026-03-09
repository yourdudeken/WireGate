package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/config"
	"github.com/yourdudeken/wiregate/internal/templates"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate deployment files",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(ConfigFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if err := cfg.Validate(); err != nil {
			fmt.Printf("Config Validation Warning: %v\n", err)
		}

		generateAction(cfg)
	},
}

type PeerContext struct {
	VPS  config.VPSConfig
	Peer config.PeerConfig
}

func generateAction(cfg *config.Config) {
	// Cleanup old deploy dir
	os.RemoveAll("deploy")

	// VPS Generation
	vpsDir := "deploy/vps"
	os.MkdirAll(filepath.Join(vpsDir, "wireguard"), 0755)
	os.MkdirAll(filepath.Join(vpsDir, "traefik_dynamic"), 0755)
	os.MkdirAll(filepath.Join(vpsDir, "letsencrypt"), 0755)

	if err := renderAndSave(vpsDir, "docker-compose.vps.yaml.tmpl", "docker-compose.yaml", cfg); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if err := renderAndSave(filepath.Join(vpsDir, "wireguard"), "wg0-server.conf.tmpl", "wg0.conf", cfg); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if err := renderAndSave(filepath.Join(vpsDir, "traefik_dynamic"), "traefik_dynamic.yaml.tmpl", "dynamic.yaml", cfg); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Peer Generation
	for _, peer := range cfg.Peers {
		peerDir := filepath.Join("deploy/peers", peer.Name)
		os.MkdirAll(filepath.Join(peerDir, "wireguard"), 0755)

		ctx := PeerContext{
			VPS:  cfg.VPS,
			Peer: peer,
		}

		if err := renderAndSave(peerDir, "docker-compose.home.yaml.tmpl", "docker-compose.yaml", ctx); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := renderAndSave(filepath.Join(peerDir, "wireguard"), "wg0-client.conf.tmpl", "wg0.conf", ctx); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	fmt.Println("Success! Deployment files generated in 'deploy/'")
}

func renderAndSave(dir, tmplName, fileName string, data interface{}) error {
	rendered, err := templates.Render(tmplName, data)
	if err != nil {
		return fmt.Errorf("rendering %s: %w", tmplName, err)
	}

	err = os.WriteFile(filepath.Join(dir, fileName), rendered, 0644)
	if err != nil {
		return fmt.Errorf("saving %s: %w", fileName, err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
