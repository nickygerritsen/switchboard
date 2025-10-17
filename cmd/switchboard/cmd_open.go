package main

import (
	"fmt"

	"github.com/nickygerritsen/switchboard/internal/logger"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <url>",
	Short: "Open a URL in the appropriate browser",
	Long: `Opens a URL by matching it against configured rules and launching
the appropriate browser with the specified profile (if any).

This is typically called by the operating system when Switchboard
is registered as the default browser.`,
	Args: cobra.ExactArgs(1),
	RunE: runOpen,
}

func runOpen(cmd *cobra.Command, args []string) error {
	url := args[0]

	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := initLogger(cfg); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer func() { _ = logger.Close() }()

	logger.Info("Opening URL: %s", url)

	rtr := createRouter(cfg)
	browserName, profile, matched := rtr.FindMatch(url)

	if !matched {
		logger.Info("No rule matched, using default browser: %s", browserName)
	}

	detector := createDetector(cfg)
	br, err := detector.Detect(browserName)
	if err != nil {
		return fmt.Errorf("failed to detect browser %q: %w", browserName, err)
	}

	lnch := createLauncher(cfg)
	if err := lnch.Launch(br, url, profile); err != nil {
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	logger.Info("Successfully opened %s in %s", url, browserName)
	return nil
}
