package server

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (server *Server) onConnect(w http.ResponseWriter, r *http.Request) {
	var uniqueId = uuid.New()
	server.logger.Printf("Received new connection attempt from %s, assiging unique id %s", r.RemoteAddr, uniqueId.String())
	server.Connect(uniqueId)

	w.WriteHeader(200)
	w.Write([]byte(uniqueId.String()))
}

func (server *Server) onDisconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		server.logger.Printf("Received disconnection attempt from %s with invalid method %s", r.RemoteAddr, r.Method)
		w.WriteHeader(200)
		return
	}

	var clientIDRaw = r.URL.Query().Get("client_id")
	if clientIDRaw == "" {
		server.logger.Printf("Received disconnection attempt from %s without client_id", r.RemoteAddr)
		w.WriteHeader(200)
		return
	}

	var clientID, err = uuid.Parse(clientIDRaw)
	if err != nil {
		server.logger.Printf("Received disconnection attempt from %s with invalid client_id %s", r.RemoteAddr, clientIDRaw)
		w.WriteHeader(200)
		return
	}
	server.Disconnect(clientID)
	w.WriteHeader(200)
}

func (server *Server) onEgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		server.logger.Printf("Received egress attempt from %s with invalid method %s", r.RemoteAddr, r.Method)
		w.WriteHeader(200)
		return
	}

	var clientIDRaw = r.URL.Query().Get("client_id")
	if clientIDRaw == "" {
		server.logger.Printf("Received egress attempt from %s without client_id", r.RemoteAddr)
		w.WriteHeader(200)
		return
	}

	var clientID, err = uuid.Parse(clientIDRaw)
	if err != nil {
		server.logger.Printf("Received egress attempt from %s with invalid client_id %s", r.RemoteAddr, clientIDRaw)
		w.WriteHeader(200)
		return
	}

	if conn, ok := server.clients[clientID.String()]; ok {
		server.clientsTimeout[clientID.String()] = time.Now()
		data, err := ioutil.ReadAll(conn.WriteBuf)
		if err != nil {
			w.WriteHeader(500)
			server.logger.Println(err)
			return
		}
		w.Write([]byte(base64.StdEncoding.EncodeToString(data)))
	} else {
		server.logger.Printf("Received egress attempt from %s with invalid client_id %s", r.RemoteAddr, clientIDRaw)
		w.WriteHeader(403)
		return
	}
}

func (server *Server) onIngress(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		server.logger.Printf("Received ingress attempt from %s with invalid method %s", r.RemoteAddr, r.Method)
		w.WriteHeader(200)
		return
	}

	var clientIDRaw = r.URL.Query().Get("client_id")
	if clientIDRaw == "" {
		server.logger.Printf("Received ingress attempt from %s without client_id", r.RemoteAddr)
		w.WriteHeader(200)
		return
	}

	var clientID, err = uuid.Parse(clientIDRaw)
	if err != nil {
		server.logger.Printf("Received ingress attempt from %s with invalid client_id %s", r.RemoteAddr, clientIDRaw)
		w.WriteHeader(200)
		return
	}

	if conn, ok := server.clients[clientID.String()]; ok {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			server.logger.Println(err)
			return
		}

		decodedData, err := base64.StdEncoding.DecodeString(string(data))
		if err != nil {
			w.WriteHeader(500)
			server.logger.Println(err)
			return
		}

		if _, err := conn.ReadBuf.Write(decodedData); err != nil {
			w.WriteHeader(500)
			server.logger.Println(err)
			return
		}
	} else {
		server.logger.Printf("Received egress attempt from %s with invalid client_id %s", r.RemoteAddr, clientIDRaw)
		w.WriteHeader(403)
		return
	}
}
