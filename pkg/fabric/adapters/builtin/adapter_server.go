package builtin

import (
	"net"

	"github.com/infinytum/go-fabric/pkg/fabric"
	"github.com/sirupsen/logrus"
)

type ServerAdapter struct {
	logger *logrus.Logger
}

func (adapter *ServerAdapter) Available() error {
	// This adapter is always ready. Awh yeah.
	return nil
}

func (adapter *ServerAdapter) Listen(network, address string) (net.Listener, error) {
	adapter.logger.Infof("Creating fabric %s listener on %s", network, address)
	return net.Listen(network, address)
}

func NewServerAdapter(logger *logrus.Logger) fabric.ServerAdapter {
	if logger == nil {
		logger = logrus.New()
	}
	logger = logger.WithField("fabric_adapter", "builtIn").WithField("fabric_mode", "server").Logger
	return &ServerAdapter{
		logger: logger,
	}
}
