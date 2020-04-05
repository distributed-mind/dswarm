// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package cli

import (
	"dswarm/dht"
	"dswarm/discovery"
	// "fmt"
	"log"
)

func run() {
	// start api
	// start discovery

	err := dht.Start()
	if err != nil {
		log.Println(err)
	}
	err = discovery.Start(
		dht.DiscoveryPacket,
		dht.DiscoveryChannel,
	)
	if err != nil {
		log.Println(err)
	}
	select {}

}
