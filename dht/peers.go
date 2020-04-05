// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package dht

import (
	"time"
	"net"
	"log"
	"bytes"
	"sort"
)

const (
	// BucketSize .
	BucketSize int = 3 // peers per bucket
)

type (
	// Peer .
	Peer struct {
		ID PeerID
		Distance Distance // TODO: distance from local
		LastSeen time.Time
		Address net.TCPAddr
	}
	sortedDistance []Distance
)

var (
	// Peers .
	Peers = make(map[PeerID]Peer)
	// ConnectedPeers .
	ConnectedPeers = make(map[PeerID]bool)

	// Bucket .
	Bucket []Peer

	// // Buckets .
	// Buckets [][BucketSize]Peer

)

// AddPeer .
func AddPeer(peer Peer) {
	if isNewPeer(peer.ID) {
		log.Println("New Peer:", peer.ID, peer.Address.Port)
	}
	Peers[peer.ID] = peer
}

func isNewPeer(id PeerID) bool {
	new := true
	for pid := range Peers {
		if bytes.Compare(pid[:], id[:]) == 0 {
			new = false
			break
		}
	}
	return new
}

// SortBucket .
func SortBucket() {
	peerDistances := make(map[Distance]PeerID)
	distances := []Distance{}
	for pid := range Peers {
		d := CalculateDistance(
			Node(ID),
			Node(pid),
		)
		peerDistances[d] = pid
		distances = append(distances, d)
	}
	// log.Println("all the distances:", distances)
	sortedPeers := []Peer{}
	sortedDistances := sortDistance(distances)
	// log.Println("sorted distances:", sortedDistances)
	for _, d := range sortedDistances {
        sortedPeers = append(sortedPeers, Peers[peerDistances[d]])
    }
	// buckets := [][BucketSize]Peer{}
	// closestPeers := [BucketSize]Peer{}
	// copy(closestPeers[:], sortedPeers[:3])
	if len(sortedPeers) >= 3 {
		Bucket = sortedPeers[:3]
	}
	// Bucket = sortedPeers[:3]
}

// SortDistance .
func sortDistance(ds []Distance) sortedDistance {
	s := sortedDistance(ds)
	sort.Sort(s)
	return s
}

func (ds sortedDistance) Len() int {
	return len(ds)
}

func (ds sortedDistance) Less(i, j int) bool {
	switch bytes.Compare(ds[i][:], ds[j][:]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		return false
	}
}

func (ds sortedDistance) Swap(i, j int) {
	ds[j], ds[i] = ds[i], ds[j]
}

// MakeBucket .
func MakeBucket() ([BucketSize]Peer) {
	//
	return [BucketSize]Peer{}
}

//
func printBucket() {
	go func() {
		status := `Three Closest Peers:
- %s
- %s
- %s
`
		for range time.Tick(5 * time.Second) {
			SortBucket()
			// log.Println(Peers)
			// log.Println(Bucket)
			if len(Bucket) >= 3 {
				log.Printf(
					status,
					Bucket[0].ID,
					Bucket[1].ID,
					Bucket[2].ID,
				)
			}	
		}
	}()
}