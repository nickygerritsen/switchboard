# Switchboard

A smart URL router that opens links in different browsers based on configurable patterns.

> **Note:** This codebase was initially created with AI assistance (Claude Code). While comprehensive tests are included and CI/CD is set up, the code has not been as thoroughly vetted as my typical projects. Use with appropriate caution and please report any issues you encounter.

## Features

- 🎯 Route URLs to different browsers based on patterns
- 🔄 URL rewriting with template variables
- 🔧 Simple YAML configuration
- 🌍 Cross-platform support (macOS, Linux, Windows)
- 👤 Profile support for Chrome based browsers and Firefox
- 🕶️ Incognito/private mode support
- 🔍 Automatic browser detection
- 📝 Debug logging support

## Installation

### Building from source

```bash
go build -o switchboard ./cmd/switchboard
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

## Shell Completion

Switchboard supports shell completion for bash, zsh, and fish. To enable completions:

### Bash

```bash
# For current session only
source <(switchboard completion bash)

# For all sessions (Linux)
switchboard completion bash > /etc/bash_completion.d/switchboard

# For all sessions (macOS)
switchboard completion bash > /usr/local/etc/bash_completion.d/switchboard
```

### Zsh

```bash
# Enable completions if not already enabled
echo "autoload -U compinit; compinit" >> ~/.zshrc

# For current session only
source <(switchboard completion zsh)

# For all sessions
switchboard completion zsh > "${fpath[1]}/_switchboard"
# Then start a new shell
```

### Fish

```bash
# For current session only
switchboard completion fish | source

# For all sessions
switchboard completion fish > ~/.config/fish/completions/switchboard.fish
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

### URL Rewriting

Rewrite URLs before opening them in browsers. This is useful for redirecting to privacy-friendly frontends or alternative services.

```yaml
rules:
  # Redirect Twitter/X to xcancel (privacy frontend)
  - match:
      - "twitter.com/*"
      - "x.com/*"
    rewrite: "xcancel.com{path}"
    browser: firefox
```

**Available template variables:**
- `{scheme}` - URL scheme (http, https)
- `{host}` - Hostname (e.g., "example.com")
- `{port}` - Port number (empty if not specified)
- `{path}` - Path portion (e.g., "/foo/bar")
- `{query}` - Query string (e.g., "key=value&key2=value2")
- `{fragment}` - Fragment/hash (e.g., "section")

**Examples:**
- `https://twitter.com/user/status/123` → `https://xcancel.com/user/status/123`
- `https://www.youtube.com/watch?v=abc` → `https://alternative.example.com/watch?v=abc`
- `https://old.example.com/some/path` → `https://new.example.com/some/path`

### Profile Support

Open URLs in specific browser profiles:

```yaml
rules:
  - match:
      - "work.company.com"
      - "*.work.company.com"
    browser: chrome
    profile: Work

  - match:
      - "personal.example.com"
    browser: firefox
    profile: Personal
```

### Incognito/Private Mode

Open URLs in incognito or private browsing mode:

```yaml
rules:
  - match:
      - "bank.com"
      - "*.bank.com"
    browser: firefox
    incognito: true
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

# Generate shell completion scripts
switchboard completion [bash|zsh|fish]
```

## License

MIT License - see LICENSE file for details
