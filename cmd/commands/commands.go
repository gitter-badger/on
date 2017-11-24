package commands

import (
	"continuul.io/adm/cmd"
	"continuul.io/adm/cmd/agent"
	"continuul.io/adm/cmd/version"
	"github.com/spf13/cobra"
)

func AddCommands(cli *cmd.Cli, cmd *cobra.Command) {
	cmd.AddCommand(
		version.NewVersionCommand(cli),
		agent.NewAgentCommand(cli),
	)
}
