package main

import (
	"continuul.io/on/cmd"
	"continuul.io/on/cmd/commands"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func setupRootCommand(rootCmd *cobra.Command) {
	rootCmd.SetFlagErrorFunc(FlagErrorFunc)
	rootCmd.SetHelpTemplate(helpTemplate)
	rootCmd.SetUsageTemplate(usageTemplate)
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

func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf(
		"docker: '%s' is not a docker command.\nSee 'docker --help'", args[0])
}

func newOnCommand(cli *cmd.Cli) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "on [OPTIONS] COMMAND [ARG...]",
		Short:            "A self-sufficient runtime for discovery",
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		Args:             noArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.ShowHelp(cmd, args)
		},
	}
	setupRootCommand(rootCmd)
	// todo: common flags go here...
	commands.AddCommands(cli, rootCmd)
	return rootCmd
}

func main() {
	cli := cmd.NewOnCli(os.Stdin, os.Stdout, os.Stderr)
	rootCmd := newOnCommand(cli)

	if err := rootCmd.Execute(); err != nil {
		if sterr, ok := err.(cmd.StatusError); ok {
			if sterr.Status != "" {
				fmt.Fprintln(os.Stderr, sterr.Status)
			}
			// StatusError should only be used for errors, and all errors should
			// have a non-zero exit status, so never exit with 0
			if sterr.StatusCode == 0 {
				os.Exit(1)
			}
			os.Exit(sterr.StatusCode)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var helpTemplate = `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

var usageTemplate = `Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  {{ .CommandPath}} [command]{{end}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}

Run "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
