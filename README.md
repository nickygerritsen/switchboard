# Switchboard

A smart URL router that opens links in different browsers based on configurable patterns.

> **Note:** This codebase was initially created with AI assistance (Claude Code). While comprehensive tests are included and CI/CD is set up, the code has not been as thoroughly vetted as my typical projects. Use with appropriate caution and please report any issues you encounter.

## Features

- 🎯 Route URLs to different browsers based on patterns
- 🔧 Simple YAML configuration
- 🌍 Cross-platform support (macOS, Linux, Windows)
- 👤 Profile support for Chrome based browsers and Firefox
- 🔍 Automatic browser detection
- 📝 Debug logging support

## Installation

### Building from source

```bash
go build -o switchboard cmd/switchboard/main.go
```

### Register as a browser

```bash
# Register Switchboard as a browser (makes it available in system settings)
switchboard register

# Remove browser registration
switchboard unregister
```

After registration, set Switchboard as your default browser:
- **macOS**: System Settings → Desktop & Dock → Default web browser
- **Linux**: Settings → Default Applications → Web Browser
- **Windows**: Settings → Apps → Default apps → Web browser

## Configuration

Switchboard looks for its configuration file at:

- **Linux**: `~/.config/switchboard/config.yaml`
- **macOS**: `~/.config/switchboard/config.yaml`
- **Windows**: `%APPDATA%\switchboard\config.yaml`

See `config.example.yaml` for a complete configuration example.

### Basic Example

```yaml
defaultBrowser: chrome

rules:
  - match:
      - "*.google.com"
      - "google.com"
    browser: firefox

  - match:
      - "github.com"
      - "*.github.com"
    browser: brave
```

## Usage

```bash
# Register as a browser
switchboard register

# Unregister as a browser
switchboard unregister

# Open a URL (typically called by the OS)
switchboard open "https://example.com"

# Test which browser would open a URL
switchboard test "https://google.com"

# List detected browsers
switchboard list-browsers

# Validate configuration
switchboard validate

# Generate example configuration
switchboard init
```

## License

MIT License - see LICENSE file for details
