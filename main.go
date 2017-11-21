package main

import (
    "continuul.io/lsr/pkg/cluster"
    "fmt"
    "log"
    "time"
    "os"
    "os/signal"
    "syscall"
    "io"
    "continuul.io/lsr/cmd/server"
    "continuul.io/lsr/cmd/commands"
)

var gracefulTimeout = 5 * time.Second
var stdOut io.Writer = os.Stdout
var stdErr io.Writer = os.Stderr

func Output(s string) error {
    _, err := stdOut.Write([]byte(s))
    return err
}

func Error(s string) error {
    _, err := stdErr.Write([]byte(s))
    return err
}

// handleSignals blocks until we get an exit-causing signal
func handleSignals(config *discovery.Config, server *discovery.Server) int {
    signalCh := make(chan os.Signal, 4)
    signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
    signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGPIPE)

    // Wait for a signal
    WAIT:
    var sig os.Signal
    select {
    case s := <-signalCh:
        sig = s
    }
    Output(fmt.Sprintf("Caught signal: %v", sig))

    // Skip SIGPIPE signals
    if sig == syscall.SIGPIPE {
        goto WAIT
    }

    // Check if we should do a graceful leave
    graceful := false
    if sig == os.Interrupt {
        graceful = true
    } else if sig == syscall.SIGTERM {
        graceful = true
    }

    // Bail fast if not doing a graceful leave
    if !graceful {
        return 1
    }

    // Attempt a graceful leave
    gracefulCh := make(chan struct{})
    Output("Gracefully shutting down agent...")
    go func() {
        if err := server.Leave(); err != nil {
            Error(fmt.Sprintf("Error: %s", err))
            return
        }
        close(gracefulCh)
    }()

    // Wait for leave or another signal
    select {
    case <-signalCh:
        return 1
    case <-time.After(gracefulTimeout):
        return 1
    case <-gracefulCh:
        return 0
    }
}

func main() {

    cmd := server.NewRootCommand()
    commands.AddCommands(cmd)

    if err := cmd.Execute(); err != nil {
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

    config := discovery.DefaultConfig()

    // Serf hard-coded configuration
    config.SerfLANConfig.MemberlistConfig.BindAddr = "127.0.0.1"
    config.SerfLANConfig.MemberlistConfig.BindPort = 8500
    config.SerfLANConfig.MemberlistConfig.SuspicionMult = 2
    config.SerfLANConfig.MemberlistConfig.ProbeTimeout = 50 * time.Millisecond
    config.SerfLANConfig.MemberlistConfig.ProbeInterval = 100 * time.Millisecond
    config.SerfLANConfig.MemberlistConfig.GossipInterval = 100 * time.Millisecond

    s, err := discovery.NewServer(config)
    if err != nil {
        log.Fatal(fmt.Errorf("Failed to start lan serf: %v", err))
    }
    defer s.Shutdown()

    handleSignals(config, s)
}
