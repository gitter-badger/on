package commands

import (
    "github.com/spf13/cobra"
    "continuul.io/lsr/cmd/version"
)

func AddCommands(cmd *cobra.Command) {
    cmd.AddCommand(
        version.NewVersionCommand(),
    )
}
