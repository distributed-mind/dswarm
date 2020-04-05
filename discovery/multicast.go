// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package discovery

import (
	//
	"net"
	"log"
	"time"
	// "strconv"
)

var (
	// Packet bytes sent via udp multicast
	Packet = []byte{}
)

func multicast(packet []byte) error {
	conn, err := net.DialUDP("udp", nil,
		&net.UDPAddr{
			IP:   []byte{239, 255, 255, 19},
			Port: 25519,
			// Zone: "udp4",
		},
	)
	if err != nil {
		log.Println(err)
		return err
	}
	err = conn.SetWriteBuffer(maxUDPSize)
	if err != nil {
		log.Println(err)
		return err
	}
	go func() {
		defer func() {
			log.Println("Stopping UDP discovery multicaster...")
			conn.Close()
		}()
		// log.Println("Started UDP discovery multicaster...")
		for range time.Tick(5 * time.Second) {
			//
			// defer conn.Close()
			err = conn.SetDeadline(time.Now().Add(1 * time.Second))
			if err != nil {
				log.Println(err)
				// return err
			}
			i, err := conn.Write(packet)
			if err != nil {
				log.Println(i, err)
				// return err
			}
			//
		}
	}()
	return nil
}
