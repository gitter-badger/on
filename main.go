package main

import (
    "fmt"
    "os"
    "continuul.io/lsr/cmd"
    "continuul.io/lsr/cmd/commands"
    "github.com/spf13/cobra"
)

func setupRootCommand(rootCmd *cobra.Command) {
    rootCmd.SetFlagErrorFunc(FlagErrorFunc)
}

// FlagErrorFunc prints an error message which matches the format of the
// docker/docker/cli error messages
func FlagErrorFunc(command *cobra.Command, err error) error {
    if err == nil {
        return nil
    }

    usage := ""
    if command.HasSubCommands() {
        usage = "\n\n" + command.UsageString()
    }
    return cmd.StatusError{
        Status:     fmt.Sprintf("%s\nSee '%s --help'.%s", err, command.CommandPath(), usage),
        StatusCode: 125,
    }
}

func newLsrCommand(cli *cmd.Cli) *cobra.Command {
    rootCmd := &cobra.Command{
        Use:              "lsr [OPTIONS] COMMAND [ARG...]",
        Short:            "A self-sufficient runtime for discovery",
        SilenceUsage:     true,
        SilenceErrors:    true,
        TraverseChildren: true,
        Args:             cmd.NoArgs,
    }
    setupRootCommand(rootCmd)
    // todo: common flags go here...
    commands.AddCommands(cli, rootCmd)
    return rootCmd
}

func main() {

    cli := cmd.NewLsrCli(os.Stdin, os.Stdout, os.Stderr)
    rootCmd := newLsrCommand(cli)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(os.Stderr, err)
        //if sterr, ok := err.(cli.StatusError); ok {
        //    if sterr.Status != "" {
        //        fmt.Fprintln(stderr, sterr.Status)
        //    }
        //    // StatusError should only be used for errors, and all errors should
        //    // have a non-zero exit status, so never exit with 0
        //    if sterr.StatusCode == 0 {
        //        os.Exit(1)
        //    }
        //    os.Exit(sterr.StatusCode)
        //}
        //fmt.Fprintln(stderr, err)
        os.Exit(1)
    }


    //dockerCli := command.NewDockerCli(stdin, stdout, stderr)
    //cmd := newDockerCommand(dockerCli)
    //
    //stdin, stdout, stderr := term.StdStreams()
}
