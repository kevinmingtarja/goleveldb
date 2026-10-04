// Package assert checks internal invariants.
package assert

import "log"

// True logs a fatal error and exits if condition is false.
func True(cond bool) {
	if !cond {
		log.Fatal("assertion failed")
	}
}
