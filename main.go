package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yuzujoe/mcp-nr-inactive-user/logger"
	"github.com/yuzujoe/mcp-nr-inactive-user/tools"
)

var (
	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "mcp-nr-inactive-user",
		Long:  `mcp-nr-inactive-user is a CLI tool to manage inactive newrelic users`,
	}

	studioCmd = &cobra.Command{
		Use:   "studio",
		Short: "Studio command",
		Long:  `Studio command is used to manage studio resources`,
		Run: func(_ *cobra.Command, _ []string) {
			logger.InitLogger()

			if err := runStudioServer(); err != nil {
				slog.Error("Error running server:", "error", err)
				return
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(studioCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}

func runStudioServer() error {
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	licenseKey := viper.GetString("NEW_RELIC_API_KEY")
	if licenseKey == "" {
		return fmt.Errorf("NEW_RELIC_API_KEY is not set")
	}

	nrClient, err := tools.NewNerdGraphServer(licenseKey)
	if err != nil {
		return fmt.Errorf("failed to create newrelic client: %w", err)
	}

	hooks := &server.Hooks{}
	mcpServer := server.NewMCPServer(
		"newrelic-mcp-server",
		"0.0.1",
		server.WithHooks(hooks),
	)

	if err := tools.ResisterTools(mcpServer, nrClient); err != nil {
		return fmt.Errorf("failed to register tools: %w", err)
	}

	if err := server.ServeStdio(mcpServer); err != nil {
		slog.Error("Error starting server:", "error", err)
		return fmt.Errorf("error starting server: %w", err)
	}

	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
