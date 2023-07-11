package cmd

import (
	"fmt"
	"os"

	cmd_server "github.com/fs714/github-star-manager/cmd/server"
	cmd_version "github.com/fs714/github-star-manager/cmd/version"
	"github.com/fs714/github-star-manager/pkg/config"
	"github.com/fs714/github-star-manager/pkg/utils/log"
	"github.com/fs714/github-star-manager/pkg/utils/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgPath string

var rootCmd = &cobra.Command{
	Use:     "github-star-manager",
	Version: version.Version,
	Short:   "A web service to manage github star in local database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	log.Sync()
}

func init() {
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false

	config.Viper = viper.New()
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")

	cmd_version.InitStartCmd()
	cmd_server.InitStartCmd()

	rootCmd.AddCommand(cmd_version.StartCmd)
	rootCmd.AddCommand(cmd_server.StartCmd)
}

func initConfig() {
	if cfgPath != "" {
		config.Viper.SetConfigFile(cfgPath)
		config.Viper.SetConfigType("yaml")
	} else {
		config.Viper.SetConfigName("github-star-manager")
		config.Viper.SetConfigType("yaml")

		dir, err := os.Getwd()
		if err != nil {
			fmt.Printf("failed to get current dir with err: %v\n", err)
			os.Exit(1)
		}

		config.Viper.AddConfigPath("/etc/github-star-manager")
		config.Viper.AddConfigPath(dir + "/conf")
	}

	err := config.Viper.ReadInConfig()
	if err != nil {
		fmt.Printf("failed to read configuration file, path: %s, err: %v\n", config.Viper.ConfigFileUsed(), err)
		os.Exit(1)
	}

	err = config.Viper.Unmarshal(&config.Config)
	if err != nil {
		fmt.Printf("failed to unmarshal config to struct, path: %s, err: %v\n", config.Viper.ConfigFileUsed(), err)
		os.Exit(1)
	}
}
