package logrotate

import (
	"fmt"
	"github.com/rickb777/logrotate/signals"
	"io"
	"log"
)

// NewLogWriter opens a log file for writing and attaches a signal handler
// (signals.RunOnPoke) to it. Therefore, when the program receives SIGHUP
// or SIGUSR1, it will close then reopen the log file, allowing log rotation
// to happen.
//
// Note that the logName may be blank or "-", in which case the defaultWriter
// will be used instead of a log file; there is no signal handler in this case.
func NewLogWriter(logName string, defaultWriter io.Writer) io.Writer {
	if logName == "" || logName == "-" {
		return defaultWriter
	}

	row := NewReopenWriter(logName)

	err := row.Open()
	if err != nil {
		panic(fmt.Errorf("Failed to open %s: %v", logName, err))
	}

	signals.RunOnPoke(
		func(message, signalName string) {
			log.Printf("%s %s", message, signalName)
		},
		func() {
			err := row.Close()
			if err != nil {
				log.Printf("Failed to close %s: %v", row.FileName(), err)
			}
			err = row.Open()
			if err != nil {
				log.Printf("Failed to open %s: %v", row.FileName(), err)
			}
		})

	return row
}
