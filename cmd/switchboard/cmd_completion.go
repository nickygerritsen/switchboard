package main

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for switchboard.

To load completions:

Bash:

  $ source <(switchboard completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ switchboard completion bash > /etc/bash_completion.d/switchboard
  # macOS:
  $ switchboard completion bash > /usr/local/etc/bash_completion.d/switchboard

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for current session only:
  $ source <(switchboard completion zsh)

  # To load completions for each session, execute once:
  $ switchboard completion zsh > "${fpath[1]}/_switchboard"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ switchboard completion fish | source

  # To load completions for each session, execute once:
  $ switchboard completion fish > ~/.config/fish/completions/switchboard.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:                  runCompletion,
}

func runCompletion(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		return rootCmd.GenBashCompletion(cmd.OutOrStdout())
	case "zsh":
		return rootCmd.GenZshCompletion(cmd.OutOrStdout())
	case "fish":
		return rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
	}
	return nil
}
