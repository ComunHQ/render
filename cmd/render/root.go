package render

import (
	"fmt"
	"os"

	"github.com/ComunHQ/render/pkg/render"
	"github.com/spf13/cobra"
)

var configFile string
var selection string

var rootCmd = &cobra.Command{
	Use:   "render",
	Short: "render - a simple CLI to render a helm template and values into yaml",
	Long:  `render - a simple CLI to render a helm template and values into yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		configs := render.GetConfigs(configFile)
		for name, config := range configs {
			if selection == "" || selection == name {
				render.Run(config)
			}
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.Flags().StringVarP(&selection, "select", "s", "", "select a specific template to render")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
