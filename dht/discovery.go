// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package dht

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"time"
)

var (
	// DiscoveryPacket .
	DiscoveryPacket = []byte{}
	// DiscoveryChannel .
	DiscoveryChannel = make(chan []byte)
)

// Endpoint .
// type Endpoint struct {
// 		ID PeerID
// 		IPs []net.IP
// 		Port int
// 	}

// EncodeDiscoveryPacket .
func EncodeDiscoveryPacket(port int) []byte {
	packet := bytes.Buffer{}
	enc := gob.NewEncoder(&packet)
	err := enc.Encode(
		struct {
			ID PeerID
			IPs []net.IP
			Port int
		}{
			ID: ID,
			Port: port,
			IPs: findPrivateIPs(),
		},
	)
	if err != nil {
		log.Println(err)
	}
	return packet.Bytes()
}

// DecodeDiscoveryPackets .
func DecodeDiscoveryPackets(c <-chan []byte) {
	go func() {
		// log.Println("Starting discovery decoder")
		for {
			b, ok := <-c
			if ok == false {
				continue
			}
			buf := bytes.NewBuffer(b)
			dec := gob.NewDecoder(buf)
			e := struct {
				ID PeerID
				IPs []net.IP
				Port int
			}{}
			err := dec.Decode(&e)
			if err != nil {
				log.Println(err)
			}
			// log.Println(e)
			if bytes.Compare(e.ID[:], ID[:]) != 0 {
				AddPeer(
					Peer{
						ID: e.ID,
						Address: net.TCPAddr{
							IP: e.IPs[0],
							Port: e.Port,
						},
						LastSeen: time.Now().UTC(),
					},
				)
			}
		}
	}()
}


func findPrivateIPs() []net.IP {
	ips := []net.IP{}
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return ips
	}
	for _, iface := range interfaces {
		ip, isPrivate := isPrivateNetwork(iface)
		if isPrivate {
			ips = append(ips, ip)
		}
	}
	return ips
}

func isPrivateNetwork(i net.Interface) (net.IP, bool) {
	// https://tools.ietf.org/html/rfc1918
	// 10.0.0.0 – 10.255.255.255
	// 10.0.0.0/8 (255.0.0.0)
	rfc1918_24 := net.IPNet{
		IP: []byte{10, 0, 0, 0},
		Mask: []byte{255, 0, 0, 0},

	}
	// 172.16.0.0 – 172.31.255.255
	// 172.16.0.0/12 (255.240.0.0)
	rfc1918_20 := net.IPNet{
		IP: []byte{172, 16, 0, 0},
		Mask: []byte{255, 240, 0, 0},

	}
	// 192.168.0.0 – 192.168.255.255
	// 192.168.0.0/16 (255.255.0.0)
	rfc1918_16 := net.IPNet{
		IP: []byte{192, 168, 0, 0},
		Mask: []byte{255, 255, 0, 0},

	}
	addrs, err := i.Addrs()
	if err != nil {
		log.Println(err)
	}
	isPrivate := false
	ip := net.IP{}
	// mask := net.IPMask{}
	for _, n := range []net.IPNet {
		rfc1918_24,
		rfc1918_20,
		rfc1918_16,
	} {
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				// mask = v.Mask
			case *net.IPAddr:
				// fail
			}
			if n.Contains(ip) {
				isPrivate = true
				break
			}
		}
		if isPrivate {
			break
		}
	}
	if isPrivate {
		// AddIP(ip, mask)
	}
	return ip, isPrivate
}
