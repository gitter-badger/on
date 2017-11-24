package discovery

// Ports is used to simplify the configuration by
// providing default ports, and allowing the addresses
// to only be specified once
type PortConfig struct {
	RPC     int // CLI RPC
	SerfLan int `mapstructure:"serf_lan"` // LAN gossip (Client + Server)
	Server  int // Server internal RPC
}

// AddressConfig is used to provide address overrides
// for specific services. By default, either ClientAddress
// or ServerAddress is used.
type AddressConfig struct {
	DNS   string // DNS Query interface
	HTTP  string // HTTP API
	HTTPS string // HTTPS API
	RPC   string // CLI RPC
}

type LANConfig struct {
	// The name of this node. This must be unique in the cluster.
	Name string

	// Configuration related to what address to bind to and ports to
	// listen on. The port is used for both UDP and TCP gossip.
	// It is assumed other nodes are running on this port, but they
	// do not need to.
	BindAddr string
	BindPort int

	// Configuration related to what address to advertise to other
	// cluster members. Used for nat traversal.
	AdvertiseAddr string
	AdvertisePort int
}
