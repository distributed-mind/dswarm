package dht

import (
	"net/url"
	"strconv"
	"crypto/tls"
	"log"
	"dswarm/id"
	"dswarm/wshs/client"
	"net/rpc/jsonrpc"
	// "crypto/rand"


	"github.com/gorilla/websocket"

)


func pingPeer(peer Peer) {
	u := &url.URL{
		Scheme: "ws",
		Host: peer.Address.IP.String() + ":" + strconv.Itoa(peer.Address.Port),
		Path: "/dht",
	}
	tc, err := tls.Dial("tcp", u.Host, &tls.Config{
		RootCAs: nil,
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println(err)
		return
	}
	ws, _, err := websocket.NewClient(tc, u, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
		return
	}
	ok := client.Shake(
		ws,
		id.Public,
		peer.ID,
		id.Private,
	)
	if ok {
		// log.Printf("%s\n", "Handshake Success: client")
	} else {
		log.Printf("%s\n", "Handshake Fail: client")
		return
	}
	c := jsonrpc.NewClient(&conn{websocket: ws})
	defer func() {
		// if err := rpc.Close() ; err != nil { // why does this hang?
		// 	log.Println(err)
		// }
		if err := ws.Close() ; err != nil {
			log.Println("ws:", err)
		}
		if err := tc.Close() ; err != nil {
			log.Println("tls:", err)
		}
		log.Println("closed client ws connection")
	}()
    reply := Message{}
    args := Message{}
    err = c.Call("PeerRPC.Ping", &args, &reply)
    if err != nil {
        log.Println("err: ", err)
    } else {
		log.Printf("Response: %x\n", reply)
	}
}

// Ping .
func (p *PeerRPC) Ping(args, r *Message) error {
	// log.Println("received ping rpc")
	// b := make([]byte, 32)
	// _, err := rand.Read(b)
	// if err != nil {
	// 	log.Println(err)
	// }
	*r = []byte("pong")
	// *r = b
    return nil
}