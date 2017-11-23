package root

import (
    "github.com/spf13/cobra"
    "continuul.io/lsr/cmd"
)

func NewRootCommand() *cobra.Command {
    //var flags *pflag.FlagSet

    command := &cobra.Command{
        Use:              "lsr [OPTIONS] COMMAND [ARG...]",
        Short:            "A self-sufficient runtime for discovery",
        SilenceUsage:     true,
        SilenceErrors:    true,
        TraverseChildren: true,
        Args:             cmd.NoArgs,
    }

    return command
}