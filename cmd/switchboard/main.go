package main

import (
	"fmt"
	"os"

	"github.com/nickygerritsen/switchboard/internal/browser"
	"github.com/nickygerritsen/switchboard/internal/config"
	"github.com/nickygerritsen/switchboard/internal/launcher"
	"github.com/nickygerritsen/switchboard/internal/logger"
	"github.com/nickygerritsen/switchboard/internal/registration"
	"github.com/nickygerritsen/switchboard/internal/router"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "switchboard",
		Short: "Smart URL router for opening links in different browsers",
		Long: `Switchboard is a smart URL router that opens links in different browsers
based on configurable patterns. It can register as your default browser
and route URLs based on domains, paths, protocols, and ports.`,
		Version: version,
	}
)

// Interfaces for dependency injection in tests
type browserDetector interface {
	Detect(name string) (*browser.Browser, error)
	DetectAll() map[string]*browser.Browser
}

type urlRouter interface {
	FindMatch(url string) (browser, profile string, incognito, matched bool, rewrittenURL string)
}

type browserLauncher interface {
	Launch(br *browser.Browser, url, profile string, incognito bool) error
}

type browserRegistrar interface {
	Register() error
	Unregister() error
	IsRegistered() (bool, error)
	GetBinaryPath() (string, error)
}

// Factory functions that can be overridden in tests
var (
	detectorFactory = func(cfg *config.Config) browserDetector {
		return browser.NewDetector(cfg)
	}
	routerFactory = func(cfg *config.Config) urlRouter {
		return router.NewRouter(cfg)
	}
	launcherFactory = func(cfg *config.Config) browserLauncher {
		return launcher.NewLauncher(cfg)
	}
	registrarFactory = func() (browserRegistrar, error) {
		return registration.New()
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is OS-specific location)")

	// Register subcommands
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(listBrowsersCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(unregisterCmd)
	rootCmd.AddCommand(completionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// loadConfig loads and validates the configuration
func loadConfig() (*config.Config, error) {
	var cfg *config.Config
	var err error

	if cfgFile != "" {
		// Load from specified file
		cfg, err = config.LoadFrom(cfgFile)
	} else {
		// Load from default location
		configPath, pathErr := config.GetConfigPath()
		if pathErr != nil {
			return nil, fmt.Errorf("failed to get config path: %w", pathErr)
		}
		cfg, err = config.LoadFrom(configPath)
	}

	if err != nil {
		return nil, err
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// initLogger initializes the logger from config
func initLogger(cfg *config.Config) error {
	return logger.Init(cfg)
}

// createDetector creates a browser detector
func createDetector(cfg *config.Config) browserDetector {
	return detectorFactory(cfg)
}

// createRouter creates a URL router
func createRouter(cfg *config.Config) urlRouter {
	return routerFactory(cfg)
}

// createLauncher creates a browser launcher
func createLauncher(cfg *config.Config) browserLauncher {
	return launcherFactory(cfg)
}
