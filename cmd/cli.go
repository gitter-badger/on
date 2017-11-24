package cmd

import (
    "io"
    "log"
    "fmt"
    "github.com/spf13/cobra"
)

// DockerCli is an instance the docker command line client.
// Instances of the client can be returned from NewDockerCli.
type Cli struct {
    out, err  io.Writer
    outLogger *log.Logger
    errLogger *log.Logger
}

func NewLsrCli(in io.Reader, out, err io.Writer) *Cli {
    logger := log.New(out, "", log.LstdFlags)
    return &Cli{
        out: out,
        err: err,
        outLogger: logger,
    }
}

func (c *Cli) Println(s string) {
    c.out.Write([]byte(fmt.Sprintf("%s\n", s)))
}

func (c *Cli) Printf(format string, a ...interface{}) {
    c.Println(fmt.Sprintf(format, a...))
}

func (c *Cli) Info(message string) {
    c.outLogger.Print("[INFO] %s", message)
}

func (c *Cli) Error(message string) {
    c.errLogger.Print("[ERROR] %s", message)
}

// ShowHelp shows the command help.
func (cli *Cli) ShowHelp(cmd *cobra.Command, args []string) error {
    cmd.SetOutput(cli.err)
    cmd.HelpFunc()(cmd, args)
    return nil
}
