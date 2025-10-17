package main

import (
	"fmt"

	"github.com/nickygerritsen/switchboard/internal/config"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long:  `Validates the configuration file and reports any errors.`,
	RunE:  runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Configuration is valid")
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  Default browser: %s\n", cfg.DefaultBrowser)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  Rules: %d\n", len(cfg.Rules))
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  Debug: %v\n", cfg.Debug)

	// Show log file location
	if cfg.LogFile != "" {
		// Resolve "auto" to actual path
		if cfg.LogFile == "auto" {
			logPath, err := config.GetLogPath()
			if err == nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  Log file: %s\n", logPath)
			} else {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  Log file: auto (failed to resolve: %v)\n", err)
			}
		} else {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  Log file: %s\n", cfg.LogFile)
		}
	} else {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  Log file: disabled")
	}

	return nil
}
