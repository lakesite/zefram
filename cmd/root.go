package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/lakesite/zefram/pkg/manager"
)

var (
	config      string
	application string

	rootCmd = &cobra.Command{
		Use:   "zefram -c [config.toml]",
		Short: "run zefram",
		Long:  `run zefram with config.toml as a daemon`,
		Run: func(cmd *cobra.Command, args []string) {
			ms := &manager.ManagerService{}
			if config == "" {
				config = "config.toml"
			}
			ms.Init(config)
			ms.Daemonize()
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of zefram",
		Long: `A number greater than 0, with prefix 'v', and possible suffixes like
            'a', 'b' or 'RELEASE'`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("zefram v0.1a")
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&config, "config", "c", "", "config file")
	rootCmd.MarkFlagRequired("config")

	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
