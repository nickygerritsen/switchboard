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

	// Detect installed browsers
	detector := createDetector(nil)
	browsers := detector.DetectAll()

	if len(browsers) == 0 {
		return fmt.Errorf("no browsers detected on your system")
	}

	// Display detected browsers
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Detected %d browser(s):\n", len(browsers))
	for name := range browsers {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", name)
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	// Ask user to choose default browser
	_, _ = fmt.Fprint(cmd.OutOrStdout(), "Enter the name of your default browser: ")
	var defaultBrowser string
	if _, err := fmt.Fscanln(cmd.InOrStdin(), &defaultBrowser); err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Validate chosen browser
	if _, ok := browsers[defaultBrowser]; !ok {
		return fmt.Errorf("browser '%s' was not detected. Please choose from the list above", defaultBrowser)
	}

	// Create config with detected browsers (using "auto" for paths)
	cfg := config.GetDefaultConfig()
	cfg.DefaultBrowser = defaultBrowser
	for name := range browsers {
		cfg.Browsers[name] = config.Browser{Path: "auto"}
	}

	// Save config
	if err := cfg.SaveTo(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nCreated configuration file at: %s\n", configPath)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Default browser set to: %s\n", defaultBrowser)
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\nEdit this file to customize your browser routing rules.")

	return nil
}
