package version

import (
    "github.com/spf13/cobra"
    "continuul.io/lsr/cmd"
    "fmt"
)

func NewVersionCommand() *cobra.Command {
    command := &cobra.Command{
        Use:   "version",
        Short: "Show version",
        Args:  cmd.NoArgs,
        Run:  func(cmd *cobra.Command, args []string) {
            showVersion()
        },
    }
    return command
}

func showVersion() {
    fmt.Printf("LSR version %s, build %s\n", Version, GitCommit)
}

