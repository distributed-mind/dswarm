// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package cli

import (
	"fmt"
)

var (
	// Version .
	Version string = "dev"
)
// PrintVersion .
func version() {
	//
	fmt.Printf("Version: %s\n", Version)
}