package agent

import (
    "github.com/spf13/cobra"
    "continuul.io/lsr/cmd"
    "fmt"
    "continuul.io/lsr/pkg/cluster"
    "time"
    "os"
    "os/signal"
    "syscall"
    "io"
    "net"
    "github.com/pkg/errors"
)

type agentOptions struct {
    NodeName     string
    BindAddr     string
    Ports        discovery.PortConfig
    StartJoin    cmd.ListOpts
    SnapshotPath string
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

func NewAgentCommand(cli *cmd.Cli) *cobra.Command {
    opts := agentOptions{
        StartJoin: cmd.NewListOpts(ValidateJoin),
    }

    command := &cobra.Command{
        Use:   "agent",
        Short: "Runs an agent",
        Args:  cmd.NoArgs,
        RunE:  func(cmd *cobra.Command, args []string) error {
            return runAgent(cli, args, &opts)
        },
    }

    flags := command.Flags()
    flags.StringVar(&opts.NodeName, "node", "",
        "node name")
    flags.StringVar(&opts.BindAddr, "bind", "",
        "address to bind server listeners to")
    flags.Var(&opts.StartJoin, "join",
        "address of agent to join on startup")
    flags.StringVar(&opts.SnapshotPath, "snapshot", "",
        "path to the snapshot file")

    return command
}

func mergeOptions(a, b *agentOptions) *agentOptions {
    var result agentOptions = *a

    if b.NodeName != "" {
        result.NodeName = b.NodeName
    }
    if b.SnapshotPath != "" {
        result.SnapshotPath = b.SnapshotPath
    }

    // Copy the bind address
    if b.BindAddr != "" {
        result.BindAddr = b.BindAddr
    }

    // Copy the start join addresses
    for _, v := range b.StartJoin.GetAllOrEmpty() {
        result.StartJoin.Set(v)
    }

    // todo: cli.Logger.Print("Join: %v", result.StartJoin.GetAll())

    return &result
}

var gracefulTimeout = 5 * time.Second
var stdOut io.Writer = os.Stdout
var stdErr io.Writer = os.Stderr

func Output(s string) error {
    fmt.Println("")
    _, err := fmt.Printf(s) //stdOut.Write([]byte(s))
    return err
}

func Info(message string) {
    Output(message)
}

func Error(s string) error {
    _, err := stdErr.Write([]byte(s))
    return err
}

// handleSignals blocks until we get an exit-causing signal
func handleSignals(_ *discovery.Config, server *discovery.Server) int {
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
    Output(fmt.Sprintf("[INFO] Caught signal: %v\n", sig))

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
    Output("Gracefully shutting down agent...\n")
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

func serverConfig(opts *agentOptions) (*discovery.Config, error) {
    base := discovery.DefaultConfig()

    base.NodeName = opts.NodeName
    base.SerfLANConfig.NodeName = base.NodeName

    if opts.BindAddr != "" {
        bindIP, bindPort, err := opts.AddrParts(opts.BindAddr)
        if err != nil {
            return nil, errors.Wrapf(err, "Invalid bind address: %s", opts.BindAddr)
        }

        base.SerfLANConfig.MemberlistConfig.BindAddr = bindIP
        base.SerfLANConfig.MemberlistConfig.BindPort = bindPort
    }

    if opts.SnapshotPath == "" {
        return nil, errors.New("Must specify data directory using --snapshot")
    }
    base.SerfLANConfig.SnapshotPath = opts.SnapshotPath

    return base, nil
}

func runAgent(cli *cmd.Cli, _ []string, opts *agentOptions) error {

    cli.Println("Starting Serf agent...")

    opts = mergeOptions(loadAgentOptions(), opts)

    config, err := serverConfig(opts)
    if err != nil {
        return fmt.Errorf("Failed to start lan serf: %s", err)
    }

    s, err := discovery.NewServer(config)
    if err != nil {
        return fmt.Errorf("Failed to start lan serf: %v", err)
    }
    defer s.Shutdown()

    s.JoinLAN(opts.StartJoin.GetAll())

    cli.Println("Serf agent running!")
    cli.Printf("     Node name: '%s'", config.NodeName)
    cli.Printf(fmt.Sprintf("     Bind addr: '%s'", opts.BindAddr))
    //Info(fmt.Sprintf("      Snapshot: %v", opts..SnapshotPath != ""))

    handleSignals(config, s)

    return nil
}

