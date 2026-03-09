package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourdudeken/wiregate/internal/config"
	"github.com/yourdudeken/wiregate/internal/wg"
)

var rotateKeysCmd = &cobra.Command{
	Use:   "rotate-keys",
	Short: "Rotate WireGuard keys",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig("config.yaml")
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		fmt.Println("Rotating VPS keys...")
		vpsKeys, _ := wg.GenerateKeyPair()
		cfg.VPS.PrivateKey = vpsKeys.Private
		cfg.VPS.PublicKey = vpsKeys.Public

		fmt.Println("Rotating peer keys...")
		for i := range cfg.Peers {
			fmt.Printf("  - Peer: %s\n", cfg.Peers[i].Name)
			keys, _ := wg.GenerateKeyPair()
			cfg.Peers[i].PrivateKey = keys.Private
			cfg.Peers[i].PublicKey = keys.Public
		}

		err = config.SaveConfig("config.yaml", cfg)
		if err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Println("Keys rotated successfully for VPS and all peers.")
		fmt.Println("Run 'generate' and 'deploy' to apply changes.")
	},
}

func init() {
	rootCmd.AddCommand(rotateKeysCmd)
}
