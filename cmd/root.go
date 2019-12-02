package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/lakesite/ls-governor"

	"github.com/lakesite/zefram/pkg/api"
	"github.com/lakesite/zefram/pkg/models"
)

var (
	config      string
	application string

	rootCmd = &cobra.Command{
		Use:   "zefram -c [config.toml]",
		Short: "run zefram",
		Long:  `run zefram with config.toml as a daemon`,
		Run: func(cmd *cobra.Command, args []string) {
			gms := &governor.ManagerService{}
			if config == "" {
				config = "config.toml"
			}

			// setup manager and create api
			gms.InitManager(config)
			gms.InitDatastore("zefram")
			gapi := gms.CreateAPI("zefram")

			// bridge logic
			model.Migrate(gapi, "zefram")
			api.SetupRoutes(gapi)

			// now daemonize the api
			gms.Daemonize(gapi)
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
