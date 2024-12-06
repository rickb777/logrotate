//go:build !linux && !darwin && !dragonfly && !freebsd && !netbsd && !openbsd
// +build !linux,!darwin,!dragonfly,!freebsd,!netbsd,!openbsd

package signals

import "log"

// Reporter provides a function to report that a signal has been received.
type Reporter func(message, signalName string)

// RunOnInterrupt only works on *Nix; otherwise it is a no-op.
func RunOnInterrupt(lgr Reporter, fn func()) {
	log.Printf("RunOnInterrupt is not implemented on this operating system\n")
}

// RunOnPoke only works on *Nix; otherwise it is a no-op.
func RunOnPoke(lgr Reporter, fn func()) {
	log.Printf("RunOnPoke is not implemented on this operating system\n")
}
