// +build linux darwin dragonfly freebsd netbsd openbsd

package signals

import (
	"os"
	"os/signal"
	"syscall"
)

// Logger provides a function to report that a signal has been received.
type Logger func(message, signalName string)

// RunOnInterrupt adds a signal listener for SIGINT, SIGTERM, SIGQUIT and SIGKILL.
// When a signal is received, fn is invoked.
// The lgr logger can be nil, which will disable logging.
func RunOnInterrupt(lgr Logger, fn func()) {
	RunOnSignal("Interrupted", lgr, fn,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
}

// RunOnPoke adds a signal listener for SIGHUP and SIGUSR1.
// When a signal is received, fn is invoked.
// The lgr logger can be nil, which will disable logging.
func RunOnPoke(lgr Logger, fn func()) {
	RunOnSignal("Poked", lgr, fn,
		syscall.SIGHUP, syscall.SIGUSR1)
}

// RunOnSignal adds a signal listener for some signals. For each signal
// received, the actionMessage is logged and the fn function is called.
// The lgr logger can be nil, which will disable logging.
func RunOnSignal(actionMessage string, lgr Logger, fn func(), signals ...os.Signal) {
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
				if lgr != nil {
					lgr(actionMessage, s.String())
				}
				fn()
			}
		}()
	}
}

// SigChannelBuffer allows the buffering to be increased. This may be necessary
// in situations where signals might be received as fast as they are processed.
// If the buffering is too small, some signals might be lost.
var SigChannelBuffer = 1
