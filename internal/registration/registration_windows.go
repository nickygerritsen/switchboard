//go:build windows

package registration

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/registry"
)

// New creates a new Registrar for the current platform
func New() (Registrar, error) {
	return newWindowsRegistrar()
}

type windowsRegistrar struct {
	binaryPath string
}

func newWindowsRegistrar() (*windowsRegistrar, error) {
	binPath, err := getBinaryPath()
	if err != nil {
		return nil, err
	}

	return &windowsRegistrar{
		binaryPath: binPath,
	}, nil
}

func (r *windowsRegistrar) Register() error {
	// Register as a browser application
	if err := r.registerApplication(); err != nil {
		return fmt.Errorf("failed to register application: %w", err)
	}

	// Register URL protocol handlers
	if err := r.registerProtocolHandlers(); err != nil {
		return fmt.Errorf("failed to register protocol handlers: %w", err)
	}

	return nil
}

func (r *windowsRegistrar) Unregister() error {
	// Remove application registration
	if err := r.unregisterApplication(); err != nil {
		return fmt.Errorf("failed to unregister application: %w", err)
	}

	// Remove protocol handlers
	if err := r.unregisterProtocolHandlers(); err != nil {
		return fmt.Errorf("failed to unregister protocol handlers: %w", err)
	}

	return nil
}

func (r *windowsRegistrar) IsRegistered() (bool, error) {
	// Check if our application key exists
	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Clients\StartMenuInternet\Switchboard`,
		registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, err
	}
	defer key.Close()

	return true, nil
}

func (r *windowsRegistrar) GetBinaryPath() (string, error) {
	return r.binaryPath, nil
}

// registerApplication registers Switchboard as a browser application
func (r *windowsRegistrar) registerApplication() error {
	// Main application key
	appKey := `Software\Clients\StartMenuInternet\Switchboard`
	key, _, err := registry.CreateKey(registry.CURRENT_USER, appKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	// Set application name
	if err := key.SetStringValue("", "Switchboard"); err != nil {
		return err
	}

	// Set capabilities key
	capKey, _, err := registry.CreateKey(registry.CURRENT_USER,
		appKey+`\Capabilities`,
		registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer capKey.Close()

	if err := capKey.SetStringValue("ApplicationName", "Switchboard"); err != nil {
		return err
	}
	if err := capKey.SetStringValue("ApplicationDescription", "Smart URL router for opening links in different browsers"); err != nil {
		return err
	}

	// Register URL associations
	urlKey, _, err := registry.CreateKey(registry.CURRENT_USER,
		appKey+`\Capabilities\URLAssociations`,
		registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer urlKey.Close()

	if err := urlKey.SetStringValue("http", "SwitchboardURL"); err != nil {
		return err
	}
	if err := urlKey.SetStringValue("https", "SwitchboardURL"); err != nil {
		return err
	}

	// Register with RegisteredApplications
	regAppsKey, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\RegisteredApplications`,
		registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer regAppsKey.Close()

	return regAppsKey.SetStringValue("Switchboard", appKey+`\Capabilities`)
}

// registerProtocolHandlers registers HTTP and HTTPS protocol handlers
func (r *windowsRegistrar) registerProtocolHandlers() error {
	protocols := []string{"http", "https"}

	for _, protocol := range protocols {
		if err := r.registerProtocol(protocol); err != nil {
			return err
		}
	}

	return nil
}

// registerProtocol registers a single protocol handler
func (r *windowsRegistrar) registerProtocol(protocol string) error {
	// Create protocol key
	protocolKey := `Software\Classes\SwitchboardURL`
	key, _, err := registry.CreateKey(registry.CURRENT_USER, protocolKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	if err := key.SetStringValue("", fmt.Sprintf("Switchboard %s Protocol", protocol)); err != nil {
		return err
	}
	if err := key.SetStringValue("URL Protocol", ""); err != nil {
		return err
	}

	// Create command key
	cmdKey, _, err := registry.CreateKey(registry.CURRENT_USER,
		protocolKey+`\shell\open\command`,
		registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer cmdKey.Close()

	cmdLine := fmt.Sprintf(`"%s" open "%%1"`, r.binaryPath)
	return cmdKey.SetStringValue("", cmdLine)
}

// unregisterApplication removes the application registration
func (r *windowsRegistrar) unregisterApplication() error {
	// Remove from RegisteredApplications
	regAppsKey, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\RegisteredApplications`,
		registry.SET_VALUE)
	if err == nil {
		_ = regAppsKey.DeleteValue("Switchboard")
		regAppsKey.Close()
	}

	// Remove application key
	return registry.DeleteKey(registry.CURRENT_USER,
		`Software\Clients\StartMenuInternet\Switchboard`)
}

// unregisterProtocolHandlers removes protocol handlers
func (r *windowsRegistrar) unregisterProtocolHandlers() error {
	return registry.DeleteKey(registry.CURRENT_USER,
		`Software\Classes\SwitchboardURL`)
}
