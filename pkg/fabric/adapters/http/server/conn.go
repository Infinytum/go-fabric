package server

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Conn struct {
	addr     Addr
	Closed   bool
	server   *Server
	ReadBuf  *bytes.Buffer
	WriteBuf *bytes.Buffer
	sync.Mutex
}

func (c *Conn) Read(b []byte) (n int, err error) {
	c.Lock()
	defer c.Unlock()
	if c.Closed {
		return 0, errors.New("Connection is closed")
	}
	return c.ReadBuf.Read(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	c.Lock()
	defer c.Unlock()
	if c.Closed {
		return 0, errors.New("Connection is closed")
	}
	return c.WriteBuf.Write(b)
}

func (c *Conn) Close() error {
	c.Lock()
	defer c.Unlock()
	c.Closed = true
	c.server.Disconnect(c.addr.id)
	return nil
}

func (c *Conn) LocalAddr() net.Addr {
	return Addr{}
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.addr
}

func (c *Conn) SetDeadline(t time.Time) error {
	c.SetReadDeadline(t)
	c.SetWriteDeadline(t)
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	c.server.logger.Println("Fabric HTTP does not support read timeouts yet")
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	c.server.logger.Println("Fabric HTTP does not support write timeouts yet")
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
