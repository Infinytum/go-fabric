package fabric

import "net"

// ServerAdapter provides methods to create listen for incoming fabric
// connections that conform to the net.Conn interface
type ServerAdapter interface {
	// Available returns whether the adapter is currently
	// in a working state or not.
	Available() error

	// Listen creates a new connection listener
	Listen(network, address string) (net.Listener, error)
}
