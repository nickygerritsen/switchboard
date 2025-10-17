package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var unregisterCmd = &cobra.Command{
	Use:   "unregister",
	Short: "Unregister Switchboard as a browser",
	Long: `Removes Switchboard's browser registration from the operating system.
This removes Switchboard from the system's browser selection and
stops it from handling HTTP/HTTPS URLs.

Note: If Switchboard is set as your default browser, you should change
your default browser before unregistering.`,
	RunE: runUnregister,
}

func runUnregister(cmd *cobra.Command, args []string) error {
	// Create platform-specific registrar
	reg, err := registrarFactory()
	if err != nil {
		return fmt.Errorf("failed to create registrar: %w", err)
	}

	// Check if registered
	isRegistered, err := reg.IsRegistered()
	if err != nil {
		return fmt.Errorf("failed to check registration status: %w", err)
	}

	if !isRegistered {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Switchboard is not currently registered as a browser.")
		return nil
	}

	// Perform unregistration
	if err := reg.Unregister(); err != nil {
		return fmt.Errorf("failed to unregister: %w", err)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Successfully unregistered Switchboard as a browser.")

	return nil
}
