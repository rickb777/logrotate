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
// The returned writer is a ReopenWriter unless the logName is blank or "-".
//
// See NewLogWriterWithSignals. If an error arises because the file cannot be
// opened, this will cause a panic.
func MustLogWriterWithSignals(logName string, defaultWriter io.Writer) io.Writer {
	w, err := NewLogWriterWithSignals(logName, defaultWriter)
	if err != nil {
		panic(fmt.Errorf("failed to open %s: %v", logName, err))
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
//
// The returned writer is a ReopenWriter unless the logName is blank or "-".
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
			Printf("%s by %s\n", message, signalName)
		},
		func() {
			err := row.Close()
			if err != nil {
				Printf("Failed to close %s: %v\n", row.FileName(), err)
			}
			err = row.Open()
			if err != nil {
				Printf("Failed to open %s: %v\n", row.FileName(), err)
			}
		})

	return row, nil
}

// Printf emits messages via log.Printf and is used by NewLogWriterWithSignals
// and MustLogWriterWithSignals, reporting whenever any signal has been received.
// Alter this before calling those functions if different behaviour is required.
var Printf = func(format string, v ...any) {
	log.Printf(format, v...)
}
