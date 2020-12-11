package signals

import (
	"os"
	"os/signal"
	"syscall"
)

// Logger provides a function to report that a signal has been received.
type Logger func(message, signalName string)

// RunOnInterrupt adds a signal listener for SIGINT, SIGTERM, SIGQUIT and SIGKILL.
func RunOnInterrupt(lgr Logger, fn func()) {
	runOnSignals("Interrupted", lgr, fn, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
}

// RunOnPoke adds a signal listener for SIGHUP and SIGUSR1.
func RunOnPoke(lgr Logger, fn func()) {
	runOnSignals("Poked", lgr, fn, syscall.SIGHUP, syscall.SIGUSR1)
}

func runOnSignals(actionMessage string, lgr Logger, fn func(), signals ...os.Signal) {
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, signals...)
		for {
			// block until there's a signal
			s := <-sigchan
			lgr(actionMessage, s.String())
			fn()
		}
	}()
}
