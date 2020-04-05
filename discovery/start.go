// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package discovery

import (
	"log"
)

const (
	maxUDPSize int = 65535
)

// Start .
func Start(packet []byte, packets chan []byte) error {
	err := listen(packets)
	if err != nil {
		log.Println(err)
		return err
	}
	Packet = packet
	err = multicast(Packet)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Started UDP multicast discovery..")
	return nil
}

