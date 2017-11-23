package agent

import (
    "github.com/spf13/cobra"
    "continuul.io/lsr/cmd"
    "fmt"
    "continuul.io/lsr/pkg/cluster"
    "time"
    "log"
    "os"
    "os/signal"
    "syscall"
    "io"
    "net"
)

type agentOptions struct {
    NodeName  string
    BindAddr  string
    Ports     discovery.PortConfig
    StartJoin cmd.ListOpts
}

// BindAddrParts returns the parts of the BindAddr that should be
// used to configure Serf.
func (opts *agentOptions) AddrParts(address string) (string, int, error) {
    checkAddr := address

    var err error

    START:
    _, _, err = net.SplitHostPort(checkAddr)
    if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
        checkAddr = fmt.Sprintf("%s:%d", checkAddr, discovery.DefaultLANSerfPort)
        goto START
    }
    if err != nil {
        return "", 0, err
    }

    // Get the address
    addr, err := net.ResolveTCPAddr("tcp", checkAddr)
    if err != nil {
        return "", 0, err
    }

    return addr.IP.String(), addr.Port, nil
}

func ValidateJoin(val string) (string, error) {
    return val, nil
}

func defaultAgentOptions() *agentOptions {
    conf := &agentOptions{
        BindAddr: "127.0.0.1",
        StartJoin: cmd.NewListOpts(ValidateJoin),
    }
    return conf
}

func loadAgentOptions() *agentOptions {
    base := defaultAgentOptions()

    if base.NodeName == "" {
        hostname, err := os.Hostname()
        if err != nil {
            fmt.Println(fmt.Sprintf("Error determining hostname: %s", err))
            return nil
        }
        base.NodeName = hostname
    }

    return base
}

func NewAgentCommand() *cobra.Command {
    opts := agentOptions{
        StartJoin: cmd.NewListOpts(ValidateJoin),
    }

    command := &cobra.Command{
        Use:   "agent",
        Short: "Runs an agent",
        Args:  cmd.NoArgs,
        RunE:  func(cmd *cobra.Command, args []string) error {
            return runAgent(args, &opts)
        },
    }

    flags := command.Flags()
    flags.StringVar(&opts.NodeName, "node", "",
        "node name")
    flags.StringVar(&opts.BindAddr, "bind", "",
        "address to bind server listeners to")
    flags.Var(&opts.StartJoin, "join",
        "address of agent to join on startup")

    return command
}

func mergeOptions(a, b *agentOptions) *agentOptions {
    var result agentOptions = *a

    if b.NodeName != "" {
        result.NodeName = b.NodeName
    }

    // Copy the bind address
    if b.BindAddr != "" {
        result.BindAddr = b.BindAddr
    }

    // Copy the start join addresses
    for _, v := range b.StartJoin.GetAllOrEmpty() {
        result.StartJoin.Set(v)
    }

    fmt.Printf("Join: %v", result.StartJoin.GetAll())

    return &result
}

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

func serverConfig(opts *agentOptions) *discovery.Config {
    base := discovery.DefaultConfig()

    base.NodeName = opts.NodeName
    base.SerfLANConfig.NodeName = base.NodeName

    if opts.BindAddr != "" {
        bindIP, bindPort, err := opts.AddrParts(opts.BindAddr)
        if err != nil {
            fmt.Println(fmt.Sprintf("Invalid bind address: %s", err))
            return nil
        }

        base.SerfLANConfig.MemberlistConfig.BindAddr = bindIP
        base.SerfLANConfig.MemberlistConfig.BindPort = bindPort
    }

    return base
}

func runAgent(args []string, opts *agentOptions) error {
    opts = mergeOptions(loadAgentOptions(), opts)
    fmt.Printf("LSR agent %s...\n", opts.BindAddr)

    config := serverConfig(opts)
    s, err := discovery.NewServer(config)
    if err != nil {
        log.Fatal(fmt.Errorf("Failed to start lan serf: %v", err))
    }
    defer s.Shutdown()

    s.JoinLAN(opts.StartJoin.GetAll())

    handleSignals(config, s)

    return nil
}

