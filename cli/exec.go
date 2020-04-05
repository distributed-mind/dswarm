// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package cli

import (
	"os"
)

// Exec .
func Exec() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
			case "version", "-version", "--version":
			{
				version()
			}
			default:
			{
				help()
			}
		}
	} else {
		run()
	}
}
