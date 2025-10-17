package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listBrowsersCmd = &cobra.Command{
	Use:   "list-browsers",
	Short: "List all detected browsers",
	Long:  `Lists all browsers that Switchboard can detect on your system.`,
	RunE:  runListBrowsers,
}

func runListBrowsers(cmd *cobra.Command, args []string) error {
	// list-browsers doesn't need config - it just detects all browsers
	detector := createDetector(nil)
	browsers := detector.DetectAll()

	if len(browsers) == 0 {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No browsers detected")
		return nil
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Detected %d browser(s):\n\n", len(browsers))
	for _, br := range browsers {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", br.Name)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "    Path: %s\n", br.Path)

		if len(br.Profiles) > 0 {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "    Profiles:\n")
			for _, profile := range br.Profiles {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "      - %s (%s)\n", profile.Name, profile.Directory)
			}
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout())
	}

	return nil
}
