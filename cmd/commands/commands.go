package commands

import (
    "github.com/spf13/cobra"
    "continuul.io/lsr/cmd/version"
    "continuul.io/lsr/cmd/agent"
)

func AddCommands(cmd *cobra.Command) {
    cmd.AddCommand(
        version.NewVersionCommand(),
        agent.NewAgentCommand(),
    )
}
