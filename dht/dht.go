// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package dht

import (
	// "crypto/rand"
	"dswarm/id"
	"encoding/base64"
	// "encoding/gob"
	// "log"
	// "net"
	"bytes"
)

const (
	// BitLength 256bit keys
	BitLength int = 256
	// KeyLength in bytes
	KeyLength int = BitLength / 8
)

type (
	// Node :
	//  - sha256 hash key
	//  - ed25519 public identity
	//  - distance metric
	//  - 32 bytes length
	Node [KeyLength]byte
	// PeerID is an ed25519 pubkey, 32 byte array
	PeerID Node
	// PeerRPC for rpc commands
	PeerRPC Node
	// Key is a sha256 hash key, 32 byte array
	Key Node
	// Distance calculation of distance between two nodes, 32 byte array
	Distance Node
	// Message .
	Message []byte
)

var (
	// ID this peer's public id
	ID PeerID = PeerID(id.Public)
)

func (p PeerID) String() string {
	return "@" + base64.StdEncoding.EncodeToString(p[:])
}

func (k Key) String() string {
	return "&" + base64.StdEncoding.EncodeToString(k[:])
}

// CalculateDistance returns distance value between two keys
func CalculateDistance(a, b Node) Distance {
	d := Distance{}
	for i := 0 ; i < KeyLength ; i++ {
        d[i] = a[i] ^ b[i]
    }
	return d
}

// GetSmaller .
func GetSmaller(a, b Node) Node {
	if bytes.Compare(a[:], b[:]) == 1 {
		return b
	}
	return a
}

// GetBigger .
func GetBigger(a, b Node) Node {
	if bytes.Compare(a[:], b[:]) == 1 {
		return a
	}
	return b
}

// IsSmallerThan is the distance shorter than d2
func (n *Node) IsSmallerThan(n2 Node) bool {
	smaller := true
	if smallest := GetSmaller(*n, n2) ; bytes.Compare(smallest[:], n2[:]) == 0 {
		smaller = false
	}
	return smaller
}

// IsBiggerThan is the distance shorter than d2
func (n *Node) IsBiggerThan(n2 Node) bool {
	bigger := true
	if biggest := GetBigger(*n, n2) ; bytes.Compare(biggest[:], n2[:]) == 0 {
		bigger = false
	}
	return bigger
}