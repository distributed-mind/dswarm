// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package dht

import (
	"dswarm/id"
	"dswarm/wshs/server"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
	// "encoding/hex"
	"fmt"

	"github.com/gorilla/websocket"
)

var (
	// Port .
	Port int
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// open port (http/websocket?)
// /dht handles dht RPC
// /api handles admin http requests

// Start the dht api
func Start() error {
	err := srv()
	if err != nil {
		log.Println(err)
		return err
	}
	DiscoveryPacket = EncodeDiscoveryPacket(Port)
	DecodeDiscoveryPackets(DiscoveryChannel)
	printBucket()
	return nil
}

func srv() error {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Println(err)
		return err
	}
	Port = l.Addr().(*net.TCPAddr).Port
	rpc.Register(&PeerRPC{})
	mux := http.NewServeMux()
	mux.HandleFunc("/dht", dhtHanlder)
	mux.HandleFunc("/api/", apiHanlder)
	https := MakeHTTPSServer(mux)
	go func() {
		defer l.Close()
		log.Println("Starting DHT on port:", strconv.Itoa(Port))
		defer log.Println("stopping https")
		for {
			err := https.ServeTLS(l, "", "")
			log.Println(err)
		}
	}()
	return nil
}

func dhtHanlder(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := ws.Close() ; err != nil {
			log.Println("ws:", err)
		}
		log.Println("closed server ws connection")
	}()
	ok := server.Shake(
		ws,
		id.Public,
		id.Private,
	)
	if ok {
		// log.Printf("%s\n", "Handshake Success: server")
		ConnectedPeers[PeerID(server.Meta.RemoteLongTermEd25519PublicKey)] = true
		defer func() {
			ConnectedPeers[PeerID(server.Meta.RemoteLongTermEd25519PublicKey)] = false
			delete(ConnectedPeers, PeerID(server.Meta.RemoteLongTermEd25519PublicKey))
		}()
	} else {
		log.Printf("%s\n", "Handshake Fail: server")	
		return
	}
	jsonrpc.ServeConn(
		&conn{websocket: ws},
	)
}


func apiHanlder(w http.ResponseWriter, r *http.Request) {
	// log.Println("request path:", r.URL.Path)
	switch r.Method {
		case http.MethodGet:
		{
			s := strings.Split(r.URL.Path, "/")
			if len(s) == 3 && s[1] == "api" && s[2] == "ping" {
				pingTopPeers()
			}
		}
		case http.MethodPut:
		{
			s := strings.Split(r.URL.Path, "/")
			if len(s) == 3 {
				log.Println("key:", s[2])
			}
			fmt.Fprintf(w, "tmp")
		}
		default:
		{
			// err
		}
	}
}

func pingTopPeers() {
	// pings all the peers in the bucket via rpc
	for _, p := range Bucket {
		go pingPeer(p)
	}
}