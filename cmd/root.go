package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var banner = `
 __      __  ________        ________         __                                           
/  \    /  |/  _____/       /  _____/  ____ _/  |_  ____ __  _  __ _____  ___.__.          
\   \/\/   /   \  ___      /   \  ____/ __ \\   __\/ __ \\ \/ \/ // __ \<   |  |          
 \        /    \    \     \    \_\  \  ___/ |  | \  ___/ \     /\  ___/ \___  |          
  \__/\  / \______  /      \______  /\___  >|__|  \___  > \/\_/  \___  >/ ____|          
       \/         \/               \/     \/           \/             \/ \/               
                                                                                           
WireGuard VPS-to-Home Gateway Tool
`

var rootCmd = &cobra.Command{
	Use:   "wiregate",
	Short: "WireGuard VPS-to-Home Gateway Tool",
	Long: banner + `
A tool that automates exposing a home server through a VPS with a public IP 
using WireGuard and Traefik.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var ConfigFile string

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", "config.yaml", "Path to config file")
}
