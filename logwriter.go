package logrotate

import (
	"fmt"
	"github.com/rickb777/logrotate/signals"
	"io"
	"log"
)

// MustLogWriterWithSignals opens a log file for writing and attaches a signal
// handler (signals.RunOnPoke) to it. Therefore, when the program receives SIGHUP
// or SIGUSR1, it will close then reopen the log file, allowing log rotation
// to happen.
//
// Note that the logName may be blank or "-", in which case the defaultWriter
// will be used instead of a log file; there is no signal handler in this case.
//
// If an error arises, this will cause a panic.
func MustLogWriterWithSignals(logName string, defaultWriter io.Writer) io.Writer {
	w, err := NewLogWriterWithSignals(logName, defaultWriter)
	if err != nil {
		panic(fmt.Errorf("Failed to open %s: %v", logName, err))
	}
	return w
}

// NewLogWriterWithSignals opens a log file for writing and attaches a signal
// handler (signals.RunOnPoke) to it. Therefore, when the program receives SIGHUP
// or SIGUSR1, it will close then reopen the log file, allowing log rotation
// to happen.
//
// Note that the logName may be blank or "-", in which case the defaultWriter
// will be used instead of a log file; there is no signal handler in this case.
func NewLogWriterWithSignals(logName string, defaultWriter io.Writer) (io.Writer, error) {
	if logName == "" || logName == "-" {
		return defaultWriter, nil
	}

	row := NewReopenWriter(logName)

	err := row.Open()
	if err != nil {
		return nil, err
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

	return row, nil
}
