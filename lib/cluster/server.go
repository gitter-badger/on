package discovery

import (
	"fmt"
	"github.com/hashicorp/serf/serf"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	serfLANSnapshot    = "serf/local.snapshot"
	DefaultLANSerfPort = 8301
)

type Config struct {
	// Node name is the name we use to advertise. Defaults to hostname.
	NodeName string

	// SerfLANConfig is the configuration for the intra-dc serf
	SerfLANConfig *serf.Config

	// LogOutput is the location to write logs to. If this is not set,
	// logs will go to stderr.
	LogOutput io.Writer
}

func DefaultConfig() *Config {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	conf := &Config{
		NodeName:      hostname,
		SerfLANConfig: serf.DefaultConfig(),
	}

	// Increase our reap interval to 3 days instead of 24h.
	conf.SerfLANConfig.ReconnectTimeout = 3 * 24 * time.Hour
	// Ensure we don't have port conflicts
	conf.SerfLANConfig.MemberlistConfig.BindPort = DefaultLANSerfPort

	conf.SerfLANConfig.MemberlistConfig.SuspicionMult = 2
	conf.SerfLANConfig.MemberlistConfig.ProbeTimeout = 50 * time.Millisecond
	conf.SerfLANConfig.MemberlistConfig.ProbeInterval = 100 * time.Millisecond
	conf.SerfLANConfig.MemberlistConfig.GossipInterval = 100 * time.Millisecond

	return conf
}

type Server struct {
	// LAN discovery event channel
	eventChLAN chan serf.Event

	// Server configuration
	config *Config

	// LAN discovery
	serfLAN *serf.Serf

	// Shutdown guards
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex

	// Logger uses the provided LogOutput
	logger *log.Logger
}

func NewServer(config *Config) (*Server, error) {
	// Ensure we have a log output and create a logger.
	if config.LogOutput == nil {
		config.LogOutput = os.Stderr
	}
	logger := log.New(config.LogOutput, "", log.LstdFlags)

	cluster := &Server{
		config:     config,
		logger:     logger,
		eventChLAN: make(chan serf.Event, 256),
		shutdownCh: make(chan struct{}),
	}

	var err error
	cluster.serfLAN, err = cluster.setupSerf(config.SerfLANConfig, cluster.eventChLAN, serfLANSnapshot)
	if err != nil {
		return nil, fmt.Errorf("Failed to start lan serf: %v", err)
	}

	go cluster.debugSerf()

	return cluster, nil
}

func (s *Server) debugSerf() {
	ticker := time.Tick(time.Second * 5)
	for {
		select {
		case <-ticker:
			log.Printf("Members: %v\n", s.serfLAN.Members())
		}
	}
}
func (s *Server) setupSerf(conf *serf.Config, ch chan serf.Event, path string) (*serf.Serf, error) {
	conf.Init()

	conf.EventCh = ch
	conf.LogOutput = s.config.LogOutput

	return serf.Create(conf)
}

func (s *Server) JoinLAN(addrs []string) (int, error) {
	count, err := s.serfLAN.Join(addrs, true)
	if err != nil {
		log.Printf("Couldn't join cluster, starting own: %v\n", err)
	}
	return count, err
}

func (s *Server) Leave() error {
	s.logger.Printf("[INFO] lsr: server starting leave")

	// Leave the LAN pool
	if s.serfLAN != nil {
		if err := s.serfLAN.Leave(); err != nil {
			s.logger.Printf("[ERR] consul: failed to leave LAN Serf cluster: %v", err)
		}
	}

	return nil
}

func (s *Server) Shutdown() error {
	s.logger.Printf("[INFO] lsr: server starting shutdown")

	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		return nil
	}

	s.shutdown = true
	close(s.shutdownCh)

	if s.serfLAN != nil {
		s.serfLAN.Shutdown()
	}

	return nil
}
