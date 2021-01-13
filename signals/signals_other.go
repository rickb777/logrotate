// +build !linux,!darwin,!dragonfly,!freebsd,!netbsd,!openbsd

package signals

// Logger provides a function to report that a signal has been received.
type Logger func(message, signalName string)

// RunOnInterrupt only works on *Nix; otherwise it is a no-op.
func RunOnInterrupt(lgr Logger, fn func()) {}

// RunOnPoke only works on *Nix; otherwise it is a no-op.
func RunOnPoke(lgr Logger, fn func()) {}
