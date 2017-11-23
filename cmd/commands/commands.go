package commands

import (
    "github.com/spf13/cobra"
    "continuul.io/lsr/cmd/version"
    "continuul.io/lsr/cmd/agent"
    "continuul.io/lsr/cmd"
)

func AddCommands(cli *cmd.Cli, cmd *cobra.Command) {
    cmd.AddCommand(
        version.NewVersionCommand(cli),
        agent.NewAgentCommand(cli),
    )
}
