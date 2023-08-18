package main

import (
	"broabroad/internal/app/config"
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"broabroad/internal/app"
)

var (
	configPath string
)

func main() {
	err := NewRootCommand().Execute()
	if err != nil {
		panic("run root command " + err.Error())
	}
}

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.NewConfig(configPath)
			if err != nil {
				return fmt.Errorf("init config: %w", err)
			}
			return app.RunMainApp(context.Background(), cfg)
		},
	}
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "config.yml", "path to config file")
	rootCmd.Version = "1.2.3"
	return rootCmd
}
