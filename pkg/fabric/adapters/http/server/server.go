package server

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Server struct {
	clients        map[string]*Conn
	clientsTimeout map[string]time.Time
	NewConnections chan *Conn
	logger         *logrus.Logger
	httpServer     *http.Server
	wg             *sync.WaitGroup
	sync.Mutex
}

func (s *Server) Accept() (net.Conn, error) {
	conn := <-s.NewConnections
	return conn, nil
}

func (s *Server) Addr() net.Addr {
	return Addr{}
}

func (s *Server) Close() error {
	if err := s.httpServer.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	s.wg.Wait()
	return nil
}

func (s *Server) Connect(id uuid.UUID) {
	s.Lock()
	defer s.Unlock()
	s.logger.Debugf("Attempting to connect client %s from server", id.String())
	client := &Conn{
		addr:     Addr{id: id},
		server:   s,
		ReadBuf:  new(bytes.Buffer),
		WriteBuf: new(bytes.Buffer),
	}
	s.clients[id.String()] = client
	s.clientsTimeout[id.String()] = time.Now()
	s.NewConnections <- client
	s.logger.Infof("Client %s has been connected from the server", id.String())
}

func (s *Server) Disconnect(id uuid.UUID) {
	s.Lock()
	defer s.Unlock()
	s.logger.Debugf("Attempting to disconnect client %s from server", id.String())
	if client, ok := s.clients[id.String()]; ok {
		if !client.Closed {
			client.Close()
		}
		delete(s.clients, id.String())
		delete(s.clientsTimeout, id.String())
	}
	s.logger.Infof("Client %s has been disconnected from the server", id.String())
}

func (s *Server) ListenAndServe(address string) *http.Server {
	s.httpServer = &http.Server{Addr: address}
	http.HandleFunc("/connect", s.onConnect)
	http.HandleFunc("/disconnect", s.onDisconnect)
	http.HandleFunc("/egress", s.onEgress)
	http.HandleFunc("/ingress", s.onIngress)

	go func() {
		for {
			clientsTimeout := make(map[string]time.Time)
			s.Lock()
			for key, value := range s.clientsTimeout {
				clientsTimeout[key] = value
			}
			s.Unlock()
			for id, t := range clientsTimeout {
				if time.Since(t) > time.Minute {
					s.logger.Warnf("Client %s has timed out and will be disconnected from the server", id)
					s.clients[id].Close()
				}
			}
			<-time.After(5 * time.Second)
		}
	}()

	go func() {
		defer s.wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return s.httpServer
}

func NewServer(logger *logrus.Logger) *Server {
	if logger == nil {
		logger = logrus.New()
	}
	logger = logger.WithField("fabric_adapter", "http").WithField("fabric_mode", "server").WithField("component", "http_server").Logger
	return &Server{
		clients:        make(map[string]*Conn),
		clientsTimeout: make(map[string]time.Time),
		logger:         logger,
		NewConnections: make(chan *Conn),
		wg:             &sync.WaitGroup{},
		Mutex:          sync.Mutex{},
	}
}
