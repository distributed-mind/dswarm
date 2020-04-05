// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package id

import (
	"encoding/base64"
)

type (
	// PublicKey ed25519 bytes
	PublicKey [32]byte
	// PrivateKey ed25519 bytes
	PrivateKey [64]byte
)

var (
	// Public id
	Public PublicKey
	// Private secret key
	Private PrivateKey
)

func init() {
	Public, Private = Generate()
}

func (k PublicKey) String() string {
	return "@" + base64.StdEncoding.EncodeToString(k[:])
}

func (k PrivateKey) String() string {
	return "^" + base64.StdEncoding.EncodeToString(k[:])
}
