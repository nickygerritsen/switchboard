package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register Switchboard as a browser",
	Long: `Registers Switchboard as a browser with the operating system.
This allows Switchboard to appear in the system's browser selection
and handle HTTP/HTTPS URLs.

After registration, you can set Switchboard as your default browser
through your system settings.`,
	RunE: runRegister,
}

func runRegister(cmd *cobra.Command, args []string) error {
	// Create platform-specific registrar
	reg, err := registrarFactory()
	if err != nil {
		return fmt.Errorf("failed to create registrar: %w", err)
	}

	// Check if already registered
	isRegistered, err := reg.IsRegistered()
	if err != nil {
		return fmt.Errorf("failed to check registration status: %w", err)
	}

	if isRegistered {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Switchboard is already registered as a browser.")
		return nil
	}

	// Perform registration
	if err := reg.Register(); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Successfully registered Switchboard as a browser!")
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\nTo set Switchboard as your default browser:")
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  • macOS: System Settings → Desktop & Dock → Default web browser")
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  • Linux: Settings → Default Applications → Web Browser")
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  • Windows: Settings → Apps → Default apps → Web browser")

	return nil
}
