// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package id

import (
	"crypto/ed25519"
	"crypto/rand"
	"log"
)

// Generate .
func Generate() (PublicKey, PrivateKey) {
	ed25519PublicKey, ed25519PrivateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	publicKey := PublicKey{}
	privateKey := PrivateKey{}
	copy(publicKey[:], ed25519PublicKey)
	copy(privateKey[:], ed25519PrivateKey)
	log.Println("Generated ID:", publicKey)
	return publicKey, privateKey
}
