# Switchboard

A smart URL router that opens links in different browsers based on configurable patterns.

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

### Register as default browser

```bash
switchboard install
```

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
