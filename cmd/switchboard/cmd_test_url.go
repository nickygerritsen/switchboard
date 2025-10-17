package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test <url>",
	Short: "Test which browser would open a URL",
	Long: `Tests which browser would be used to open a URL without actually launching it.
Shows which rule matched and which browser/profile would be used.`,
	Args: cobra.ExactArgs(1),
	RunE: runTest,
}

func runTest(cmd *cobra.Command, args []string) error {
	url := args[0]

	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	rtr := createRouter(cfg)
	browserName, profile, incognito, matched, rewrittenURL := rtr.FindMatch(url)

	if matched {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "URL: %s\n", url)
		if rewrittenURL != url {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Rewritten URL: %s\n", rewrittenURL)
		}
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Browser: %s\n", browserName)
		if profile != "" {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Profile: %s\n", profile)
		}
		if incognito {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Incognito: yes")
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Matched: yes")
	} else {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "URL: %s\n", url)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Browser: %s (default)\n", browserName)
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Matched: no")
	}

	return nil
}
