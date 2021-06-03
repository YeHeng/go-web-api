package main

import (
	"fmt"
	"os"

	"github.com/YeHeng/go-web-api/cmd/mfmt"
	"github.com/YeHeng/go-web-api/cmd/mysqlmd"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "go-web-api",
		Short: "A generator for gin web api",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-web-api.yaml)")
	rootCmd.AddCommand(mysqlmd.NewCmdMySQLMd())
	rootCmd.AddCommand(mfmt.NewFmtCmd())
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		// Search config in home directory with name ".go-web-api" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-web-api")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
