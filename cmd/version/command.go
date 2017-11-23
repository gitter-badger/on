package version

import (
    "github.com/spf13/cobra"
    "continuul.io/lsr/cmd"
)

func NewVersionCommand(cli *cmd.Cli) *cobra.Command {
    command := &cobra.Command{
        Use:   "version",
        Short: "Show version",
        Args:  cmd.NoArgs,
        Run:  func(cmd *cobra.Command, args []string) {
            showVersion(cli)
        },
    }
    return command
}

func showVersion(cli *cmd.Cli) {
    cli.Printf("LSR version %s, build %s\n", Version, GitCommit)
}

