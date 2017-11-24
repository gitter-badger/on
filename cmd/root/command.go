package root

import (
	"continuul.io/adm/cmd"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:              "adm [OPTIONS] COMMAND [ARG...]",
		Short:            "A self-sufficient runtime for discovery",
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		Args:             cmd.NoArgs,
	}
	return command
}
