package builtin

import (
	"log"
	"net"
	"os"

	"github.com/infinytum/go-fabric/pkg/fabric"
)

type ClientAdapter struct {
	logger *log.Logger
}

func (adapter *ClientAdapter) Available() error {
	// This adapter is always ready. Awh yeah.
	return nil
}

func (adapter *ClientAdapter) Dial(network, address string) (net.Conn, error) {
	adapter.logger.Printf("Creating fabric %s connection to %s", network, address)
	return net.Dial(network, address)
}

func NewClientAdapter(logger *log.Logger) fabric.ClientAdapter {
	if logger == nil {
		logger = log.New(os.Stdout, "[BuiltInClientAdapter] ", log.Default().Flags())
	}
	logger.SetPrefix("[BuiltInClientAdapter] ")
	return &ClientAdapter{
		logger: logger,
	}
}
