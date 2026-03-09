package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var hubCmd = &cobra.Command{
	Use:   "hub",
	Short: "Manage multiple gateway contexts (hubs)",
}

var listHubsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available gateway configurations",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := os.ReadDir(".")
		if err != nil {
			fmt.Printf("Error reading directory: %v\n", err)
			return
		}

		fmt.Println("Available Gateway Hubs:")
		found := false
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") && file.Name() != "go.mod" && file.Name() != "go.sum" {
				name := strings.TrimSuffix(file.Name(), ".yaml")
				current := ""
				if file.Name() == ConfigFile {
					current = " (active)"
				}
				fmt.Printf("  - %s [%s]%s\n", name, file.Name(), current)
				found = true
			}
		}

		if !found {
			fmt.Println("No configuration files found. Use 'init' to create one.")
		}

		fmt.Println("\nTo use a specific hub, use the -c flag:")
		fmt.Println("  wiregate -c <filename>.yaml [command]")
	},
}

func init() {
	hubCmd.AddCommand(listHubsCmd)
	rootCmd.AddCommand(hubCmd)
}
