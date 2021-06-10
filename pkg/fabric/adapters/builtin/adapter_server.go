package builtin

import (
	"log"
	"net"
	"os"

	"github.com/infinytum/go-fabric/pkg/fabric"
)

type ServerAdapter struct {
	logger *log.Logger
}

func (adapter *ServerAdapter) Available() error {
	// This adapter is always ready. Awh yeah.
	return nil
}

func (adapter *ServerAdapter) Listen(network, address string) (net.Listener, error) {
	adapter.logger.Printf("Creating fabric %s listener on %s", network, address)
	return net.Listen(network, address)
}

func NewServerAdapter(logger *log.Logger) fabric.ServerAdapter {
	if logger == nil {
		logger = log.New(os.Stdout, "[HTTPServerAdapter] ", log.Default().Flags())
	}
	logger.SetPrefix("[HTTPServerAdapter] ")
	return &ServerAdapter{
		logger: logger,
	}
}
