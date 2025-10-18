# Claude Development Guide for Switchboard

This document outlines the development workflow, conventions, and best practices for working on the Switchboard project.

## Project Overview

Switchboard is a cross-platform URL router written in Go 1.24.0 that routes URLs to different browsers based on configurable rules. It supports browser profiles, incognito mode, and URL rewriting.

## GitHub Workflow

We follow a standard issue-based workflow:

1. **Start with an issue**: Always begin work by referencing a GitHub issue number
2. **Create a feature branch**: Branch from `main` using the naming convention below
3. **Implement the feature**: Make changes, write tests, ensure linting passes
4. **Create a Pull Request**: Open a PR to merge back into `main`
5. **Merge and cleanup**: After merge, update `main` and delete stale branches

### Branch Naming Conventions

- **Feature branches**: `feature/<descriptive-name>` (e.g., `feature/url-rewriting`)
- **Documentation branches**: `docs/<descriptive-name>` (e.g., `docs/url-rewriting-documentation`)
- **Bug fix branches**: `fix/<descriptive-name>` or `bugfix/<descriptive-name>`

The main branch is `main`.

## Testing Requirements

- **All code must have tests**: New features require comprehensive test coverage
- **Tests must pass**: Run `go test ./...` before committing
- **Update existing tests**: When changing function signatures, update all affected tests
- **Test files**: Place tests alongside the code they test (e.g., `router.go` → `router_test.go`)

## Linting Requirements

- **golangci-lint must pass**: Run `golangci-lint run` before committing
- **No warnings allowed**: All linter warnings must be addressed
- **Fix issues promptly**: If the linter finds issues, fix them before proceeding

## Build Requirements

- **Code must build**: Run `go build ./...` to verify all packages compile
- **No build errors**: Address all compilation errors before committing

## Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/) format for better changelog generation:

```
<type>: <description>

[optional body]

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

### Commit Types

- **feat**: New feature (appears in changelog under "New Features")
- **fix**: Bug fix (appears in changelog under "Bug Fixes")
- **docs**: Documentation changes (appears in changelog under "Documentation")
- **refactor**: Code refactoring without behavior change (appears in changelog under "Enhancements")
- **perf**: Performance improvements (appears in changelog under "Enhancements")
- **test**: Adding or updating tests (excluded from changelog)
- **chore**: Maintenance tasks (excluded from changelog)

### Examples

```
feat: add shell completion support for bash, zsh, and fish

fix: resolve Safari incognito warning not being logged

docs: update README with installation instructions

refactor: simplify router matching logic
```

Keep commit messages concise and descriptive.

## Pull Request Process

1. **Create PR from feature branch**: Use descriptive titles
2. **Include context**: Reference the issue number in the PR description
3. **Wait for review**: Don't merge immediately (unless instructed)
4. **Merge when approved**: Typically the user will merge PRs
5. **Clean up branches**: After merge, delete both local and remote feature branches

## Development Best Practices

### Tool Preferences

- **Use Edit tool for file modifications**: Don't use bash commands like `sed`, `awk`, or `echo >>`
- **Use Read tool for reading files**: Don't use bash commands like `cat`, `head`, or `tail`
- **Use Write tool for new files**: Don't use bash commands with heredocs
- **Use Bash only for actual terminal operations**: git commands, running tests, building, etc.

### Code Organization

- **Internal packages**: Core logic goes in `internal/` (e.g., `internal/router`, `internal/rewriter`)
- **Command handlers**: CLI commands live in `cmd/switchboard/`
- **Config handling**: Configuration logic in `internal/config/`
- **Keep packages focused**: Each package should have a clear, single responsibility

### When Changing Interfaces

If you modify a function signature that's part of an interface:

1. Update the interface definition (e.g., in `main.go`)
2. Update all implementations (e.g., the actual struct method)
3. Update all test doubles/fakes (e.g., in `testhelpers_test.go`)
4. Update all call sites throughout the codebase
5. Update all tests that call the function

### Documentation

- **Update README.md**: Document new features with examples
- **Update config.example.yaml**: Add example configurations for new features
- **Code comments**: Use clear, concise comments for complex logic
- **No unnecessary docs**: Don't create extra documentation files unless explicitly requested

## Common Patterns

### Router Pattern

The router uses a first-match-wins pattern:
- Rules are evaluated in order
- First matching rule wins
- If no rule matches, use the default browser

### Configuration Loading

- Config file is YAML-based
- Lives in `~/.config/switchboard/config.yaml`
- Example config in `config.example.yaml`

### Logging

- Use the `internal/logger` package
- Log levels: Debug, Info, Warn, Error
- Configure via config file or command-line flags

## Project Structure

```
switchboard/
├── cmd/
│   └── switchboard/        # CLI commands and main entry point
├── internal/
│   ├── browser/           # Browser launching logic
│   ├── config/            # Configuration handling
│   ├── logger/            # Logging utilities
│   ├── matcher/           # URL pattern matching
│   ├── rewriter/          # URL rewriting with templates
│   └── router/            # Rule matching and routing
├── README.md              # User documentation
├── config.example.yaml    # Example configuration
└── claude.md             # This file
```

## Quick Reference Commands

```bash
# Run tests
go test ./...

# Run linter
golangci-lint run

# Build all packages
go build ./...

# Build the binary
go build -o switchboard ./cmd/switchboard

# Run the application
./switchboard open <url>

# Test a URL without opening
./switchboard test-url <url>

# Validate configuration
./switchboard validate
```

## Working with Issues

1. Fetch issue details using the GitHub MCP tool
2. Understand requirements fully before starting
3. Create a feature branch
4. Implement with tests
5. Create PR when complete
6. Reference issue number in PR

## After PR Merge

The typical cleanup process:

```bash
git checkout main
git pull
git branch -d <feature-branch>
git push origin --delete <feature-branch>
```

---

This guide should help maintain consistency across development sessions. Update it as new patterns emerge or workflows change.
