package commands

import (
	"continuul.io/lsr/cmd"
	"continuul.io/lsr/cmd/agent"
	"continuul.io/lsr/cmd/version"
	"github.com/spf13/cobra"
)

func AddCommands(cli *cmd.Cli, cmd *cobra.Command) {
	cmd.AddCommand(
		version.NewVersionCommand(cli),
		agent.NewAgentCommand(cli),
	)
}
