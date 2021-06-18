package http

import (
	"net"

	"github.com/infinytum/go-fabric/pkg/fabric"
	"github.com/infinytum/go-fabric/pkg/fabric/adapters/http/server"
	"github.com/sirupsen/logrus"
)

type BuiltInServerAdapter struct {
	logger *logrus.Logger
}

func (adapter *BuiltInServerAdapter) Available() error {
	// This adapter is always ready. Awh yeah.
	return nil
}

func (adapter *BuiltInServerAdapter) Listen(network, address string) (net.Listener, error) {
	adapter.logger.Infof("Creating fabric %s listener on %s", network, address)
	server := server.NewServer(nil)
	go server.ListenAndServe(address)
	return server, nil
}

func NewServerAdapter(logger *logrus.Logger) fabric.ServerAdapter {
	if logger == nil {
		logger = logrus.New()
	}
	logger = logger.WithField("fabric_adapter", "http").WithField("fabric_mode", "server").Logger
	return &BuiltInServerAdapter{
		logger: logger,
	}
}
