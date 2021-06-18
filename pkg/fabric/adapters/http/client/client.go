package client

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Conn struct {
	addr     Addr
	endpoint string
	logger   *logrus.Logger
}

func (c *Conn) Connect() error {
	c.logger.Debug("Connecting to fabric http server")
	resp, err := http.Get(c.endpoint + "/connect")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if id, err := uuid.ParseBytes(body); err == nil {
		c.addr = Addr{id: id}

		c.logger.Info("Connected to fabric http server")
		return nil
	}
	return errors.New("invalid UUID")
}

func (c *Conn) Read(b []byte) (n int, err error) {
	resp, err := http.Get(c.endpoint + "/egress?client_id=" + url.QueryEscape(c.addr.id.String()))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	data, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, io.EOF
	}
	return copy(b, data), nil
}

func (c *Conn) Write(b []byte) (int, error) {
	var encoded = base64.StdEncoding.EncodeToString(b)
	_, err := http.Post(c.endpoint+"/ingress?client_id="+c.addr.id.String(), "text", bytes.NewBufferString(encoded))
	return len(b), err
}

func (c *Conn) Close() error {
	_, err := http.Post(c.endpoint+"/disconnect?client_id="+c.addr.id.String(), "text", new(bytes.Buffer))
	return err
}

func (c *Conn) LocalAddr() net.Addr {
	return c.addr
}

func (c *Conn) RemoteAddr() net.Addr {
	return Addr{}
}

func (c *Conn) SetDeadline(t time.Time) error {
	c.SetReadDeadline(t)
	c.SetWriteDeadline(t)
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	c.logger.Warn("Fabric HTTP does not support read timeouts yet")
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	c.logger.Warn("Fabric HTTP does not support write timeouts yet")
	return nil
}

type Addr struct {
	id uuid.UUID
}

func (a Addr) Network() string {
	return "fabric"
}

func (a Addr) String() string {
	return a.id.String()
}

func NewConn(endpoint string, logger *logrus.Logger) *Conn {
	if logger == nil {
		logger = logrus.New()
	}
	logger = logger.WithField("fabric_adapter", "http").WithField("fabric_mode", "client").WithField("component", "http_conn").Logger
	return &Conn{
		addr:     Addr{id: uuid.Nil},
		endpoint: endpoint,
		logger:   logger,
	}
}
