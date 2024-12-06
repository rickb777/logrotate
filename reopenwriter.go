package logrotate

import (
	"fmt"
	"github.com/rickb777/logrotate/safe"
	"io"
	"os"
)

// ReopenWriter is a WriteCloser that is able to close and reopen cleanly
// whilst in use. It is intended for writing log files, so always appends
// to existing files instead of truncating them.
//
// This will work well with the 'logrotate' utility on Linux. On rotation,
// 'logrotate' should rename the old file and then signal to the application,
// usually via SIGHUP or SIGUSR1.
type ReopenWriter interface {
	io.WriteCloser
	io.StringWriter
	FileName() string
	Open() error
}

type reopener struct {
	fileName string
	writer   *safe.Safe
}

// NewReopenWriter returns a new ReopenWriter with a given filename.
// The filename stays constant even when the file is closed
// and reopened.
//
// This function does not register any signal handling. For
// writing logfiles with Unix logrotate, see NewLogWriterWithSignals.
func NewReopenWriter(fileName string) ReopenWriter {
	return &reopener{fileName, safe.New(nil)}
}

// FileName returns the file name.
func (r *reopener) FileName() string {
	return r.fileName
}

// WriteString satisfies the io.StringWriter interface.
// This will block until the file is open.
func (r *reopener) WriteString(s string) (n int, err error) {
	return r.Write([]byte(s))
}

// Write satisfies the io.Writer interface.
// This will block until the file is open.
func (r *reopener) Write(p []byte) (n int, err error) {
	w := r.writer.GetWhenDefined().(io.Writer)
	return w.Write(p)
}

// Open opens the file.
func (r *reopener) Open() error {
	flag := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	w, err := os.OpenFile(r.fileName, flag, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", r.fileName, err)
	}
	r.writer.Put(w)
	return nil
}

// Close closes the file. Subsequent Write operations are paused until the
// file is reopened.
func (r *reopener) Close() error {
	w := r.writer.Get().(io.Closer)
	if w == nil {
		return fmt.Errorf("attempt to close %s when it is not open", r.fileName)
	}
	r.writer.Put(nil)
	err := w.Close()
	if err != nil {
		return fmt.Errorf("failed to close %s: %w", r.fileName, err)
	}
	return nil
}
