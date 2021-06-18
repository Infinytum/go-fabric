package http

import (
	"net"

	"github.com/infinytum/go-fabric/pkg/fabric"
	"github.com/infinytum/go-fabric/pkg/fabric/adapters/http/client"
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
	client := client.NewConn(network+"://"+address, nil)
	if err := client.Connect(); err != nil {
		return nil, err
	}
	return client, nil
}

func NewClientAdapter(logger *logrus.Logger) fabric.ClientAdapter {
	if logger == nil {
		logger = logrus.New()
	}
	logger = logger.WithField("fabric_adapter", "http").WithField("fabric_mode", "client").Logger
	return &ClientAdapter{
		logger: logger,
	}
}
