package fabric

import "net"

// ClientAdapter provides methods to create new connections that
// conform to the net.Conn interface
type ClientAdapter interface {
	// Available returns whether the adapter is currently
	// in a working state or not.
	Available() error

	// Dial creates a new connection
	Dial(network, address string) (net.Conn, error)
}
