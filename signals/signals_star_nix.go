//go:build linux || darwin || dragonfly || freebsd || netbsd || openbsd
// +build linux darwin dragonfly freebsd netbsd openbsd

package signals

import (
	"os"
	"os/signal"
	"syscall"
)

// Reporter provides a function to report that a signal has been received.
type Reporter func(message, signalName string)

// RunOnInterrupt adds a signal listener for SIGINT, SIGTERM, SIGQUIT and SIGKILL.
// When a signal is received, fn is invoked.
//
// If the rep logger is not nil, it is used to report every signal received.
func RunOnInterrupt(rep Reporter, fn func()) {
	RunOnSignal("Interrupted", rep, fn,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
}

// RunOnPoke adds a signal listener for SIGHUP and SIGUSR1.
// When a signal is received, fn is invoked.
//
// If the rep logger is not nil, it is used to report every signal received.
func RunOnPoke(rep Reporter, fn func()) {
	RunOnSignal("Poked", rep, fn,
		syscall.SIGHUP, syscall.SIGUSR1)
}

// RunOnSignal adds a signal listener for some signals. For each signal
// received, the actionMessage is logged and the fn function is called.
// The rep logger can be nil, which will disable reporting of signals.
func RunOnSignal(actionMessage string, rep Reporter, fn func(), signals ...os.Signal) {
	if len(signals) > 0 {
		if SigChannelBuffer < 1 {
			panic("The channel must be buffered")
		}
		go func() {
			sigchan := make(chan os.Signal, SigChannelBuffer*len(signals))
			signal.Notify(sigchan, signals...)
			for {
				// block until there's a signal
				s := <-sigchan
				if rep != nil {
					rep(actionMessage, s.String())
				}
				fn()
			}
		}()
	}
}

// SigChannelBuffer allows the buffering to be increased. This may be necessary
// in situations where signals might be received as fast as they are processed.
// If the buffering is too small, some signals might be lost. This is 1 by
// default, which is the minimum necessary.
var SigChannelBuffer = 1
