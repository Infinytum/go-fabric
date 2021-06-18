package builtin

import (
	"net"

	"github.com/infinytum/go-fabric/pkg/fabric"
	"github.com/sirupsen/logrus"
)

type ClientAdapter struct {
	logger *logrus.Logger
}

func (adapter *ClientAdapter) Available() error {
	// This adapter is always ready. Awh yeah.
	return nil
}

func (adapter *ClientAdapter) Dial(network, address string) (net.Conn, error) {
	adapter.logger.Infof("Creating fabric %s connection to %s", network, address)
	return net.Dial(network, address)
}

func NewClientAdapter(logger *logrus.Logger) fabric.ClientAdapter {
	if logger == nil {
		logger = logrus.New()
	}
	logger = logger.WithField("fabric_adapter", "builtIn").WithField("fabric_mode", "client").Logger
	return &ClientAdapter{
		logger: logger,
	}
}
