package http

import (
	"log"
	"net"
	"os"

	"github.com/infinytum/go-fabric/pkg/fabric"
	"github.com/infinytum/go-fabric/pkg/fabric/adapters/http/client"
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
	client := client.NewConn(network+"://"+address, nil)
	if err := client.Connect(); err != nil {
		return nil, err
	}
	return client, nil
}

func NewClientAdapter(logger *log.Logger) fabric.ClientAdapter {
	if logger == nil {
		logger = log.New(os.Stdout, "[HTTPClientAdapter] ", log.Default().Flags())
	}
	logger.SetPrefix("[HTTPClientAdapter] ")
	return &ClientAdapter{
		logger: logger,
	}
}
