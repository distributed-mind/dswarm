// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package discovery

import (
	// "bytes"
	// "encoding/base64"
	// "encoding/gob"
	"log"
	"net"
	// "strconv"
	// "strings"
)

func listen(packets chan<- []byte) error {
	conn, err := net.ListenMulticastUDP("udp", nil, 
		&net.UDPAddr{
			IP: []byte{239, 255, 255, 19},
			Port: 25519,
			// Zone: "udp4",
		},
	)
	if err != nil {
		log.Println(err)
		return err
	}
	err = conn.SetReadBuffer(maxUDPSize)
	if err != nil {
		log.Println(err)
		return err
	}
	go func() {
		defer func() {
			log.Println("Stopping UDP discovery multicast listener...")
			conn.Close()
		}()
		// log.Println("Started UDP discovery multicast listener...")
		for {
			b := make([]byte, maxUDPSize)
			l, _, err := conn.ReadFromUDP(b)
			// _, _, err := conn.ReadFromUDP(b)
			if err != nil {
				log.Println(err)
				continue
			} else {
				go func() {
					packets <- b[:l]
					// log.Println("Seen Packet:", address.String(), l)
					// buf := bytes.NewBuffer(b[:l])
					// dec := gob.NewDecoder(buf)
					// err = dec.Decode(&q)
					// if err != nil {
					// 	log.Fatal("decode error:", err)
					// }
					// s := strings.Split(string(b[:l]), "@")
					// if len(s) == 2 {
					// 	// port, err := strconv.Atoi(s[0])
					// 	_, err := strconv.Atoi(s[0])
					// 	if err != nil {
					// 		log.Println(err)
					// 		return
					// 	}
					// 	k, err := base64.StdEncoding.DecodeString(
					// 		strings.Split(s[1], ".")[0],
					// 	)
					// 	if err != nil {
					// 		log.Println(err)
					// 		return
					// 	}
					// 	if len(k) == 32 {
					// 		id := [32]byte{}
					// 		copy(id[:], k)
					// 		// if _, ok := SeenPeers[id] ; !ok && isNewPeer(id) {
					// 		// 	if id != dht.Me.PublicKey {
					// 		// 		log.Println("Found new peer:", "@" + s[1])
					// 		// 		// SeenPeers[id] = Peer{
					// 		// 		// 	ID: id,
					// 		// 		// 	Listeners: make(map[listenerType]listener),
					// 		// 		// }
					// 		// 		// SeenPeers[id].Listeners[webSocketListener] = listener{
					// 		// 		// 	address: address.IP,
					// 		// 		// 	port: port,
					// 		// 		// }
					// 		// 	}
					// 		// }
					// 	}
					// }
				}()
			}
		}
	}()
	return nil
}

