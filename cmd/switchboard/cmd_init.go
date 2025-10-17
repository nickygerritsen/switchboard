package main

import (
	"fmt"
	"os"

	"github.com/nickygerritsen/switchboard/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create example configuration file",
	Long:  `Creates an example configuration file at the default location.`,
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("configuration file already exists at %s", configPath)
	}

	// Get default config
	cfg := config.GetDefaultConfig()

	// Save config
	if err := cfg.SaveTo(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Created configuration file at: %s\n", configPath)
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\nEdit this file to customize your browser routing rules.")

	return nil
}
