package commands

import (
	"continuul.io/on/cmd"
	"continuul.io/on/cmd/agent"
	"continuul.io/on/cmd/version"
	"github.com/spf13/cobra"
)

func AddCommands(cli *cmd.Cli, cmd *cobra.Command) {
	cmd.AddCommand(
		version.NewVersionCommand(cli),
		agent.NewAgentCommand(cli),
	)
}
