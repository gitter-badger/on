package main

import (
    "fmt"
    "os"
    "continuul.io/lsr/cmd/commands"
    "continuul.io/lsr/cmd/root"
)

func main() {

    rootCmd := root.NewRootCommand()
    commands.AddCommands(rootCmd)

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
