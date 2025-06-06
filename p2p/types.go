package p2p

import (
	"strings"

	maddr "github.com/multiformats/go-multiaddr"
)

var _ P2PConfig = &Config{}

// A new type we need for writing a custom flag parser
type addrList []maddr.Multiaddr

// Config is configuration for P2P
type Config struct {
	RendezvousString string
	Port             int
	BootstrapPeers   addrList
	ExternalIP       string
}

func (c *Config) GetRendezvous() string {
	return c.RendezvousString
}

func (c *Config) GetP2PPort() int {
	return c.Port
}

func (c *Config) GetBootstrapPeers() ([]maddr.Multiaddr, error) {
	return c.BootstrapPeers, nil
}

func (c *Config) GetExternalIP() string {
	return c.ExternalIP
}

// String implement fmt.Stringer
func (al *addrList) String() string {
	addresses := make([]string, len(*al))
	for i, addr := range *al {
		addresses[i] = addr.String()
	}
	return strings.Join(addresses, ",")
}

// Set add the given value to addList
func (al *addrList) Set(value string) error {
	addr, err := maddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}
