package http

import (
	"log"
	"net"
	"os"

	"github.com/infinytum/go-fabric/pkg/fabric"
	"github.com/infinytum/go-fabric/pkg/fabric/adapters/http/server"
)

type BuiltInServerAdapter struct {
	logger *log.Logger
}

func (adapter *BuiltInServerAdapter) Available() error {
	// This adapter is always ready. Awh yeah.
	return nil
}

func (adapter *BuiltInServerAdapter) Listen(network, address string) (net.Listener, error) {
	adapter.logger.Printf("Creating fabric %s listener on %s", network, address)
	server := server.NewServer(nil)
	go server.ListenAndServe(address)
	return server, nil
}

func NewServerAdapter(logger *log.Logger) fabric.ServerAdapter {
	if logger == nil {
		logger = log.New(os.Stdout, "[HTTPServerAdapter] ", log.Default().Flags())
	}
	logger.SetPrefix("[HTTPServerAdapter] ")
	return &BuiltInServerAdapter{
		logger: logger,
	}
}
