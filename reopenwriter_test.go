package logrotate

import (
	"fmt"
	. "github.com/onsi/gomega"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

const line1 = "So shaken as we are,\n"
const line2 = "So wan with care.\n"

func TestReopenWriter(t *testing.T) {
	g := NewGomegaWithT(t)

	ch := make(chan bool)
	defer os.Remove("test.txt")
	rw := NewReopenWriter("test.txt")

	g.Expect(rw.FileName()).Should(Equal("test.txt"))

	err := rw.Open()
	g.Expect(err).Should(BeNil())

	n, err := rw.WriteString(line1)
	g.Expect(err).Should(BeNil())
	g.Expect(n).Should(Equal(len(line1)))

	go func() {
		<-ch
		rw.Close()
		time.Sleep(time.Millisecond)
		rw.Open()
		<-ch
	}()

	ch <- true

	n, err = rw.WriteString(line2)
	g.Expect(err).Should(BeNil())

	ch <- true

	rw.Close()

	file, err := os.Open("test.txt")
	g.Expect(err).Should(BeNil())

	content, err := io.ReadAll(file)
	g.Expect(err).Should(BeNil())

	g.Expect(file.Close()).Should(BeNil())

	g.Expect(string(content)).Should(Equal(line1 + line2))
}

func ExampleMustLogWriterWithSignals() {
	// This shows using MustLogWriterWithSignals to obtain an io.Writer
	// that is used for log or log/slog logging.

	defer os.Remove("example.log") // this is just to keep the example tidy

	// Open the log writer and register it to handle SIGHUP & SIGUSR1.
	// Note that os.Stdout isn't actually used in this example but might be
	// needed if the filename is a configuration parameter that might take
	// the special value "-" to indicate stdout.
	lw := MustLogWriterWithSignals("example.log", os.Stdout)

	// lgr is the 'basic' log package
	lgr := log.New(lw, "", 0)

	// we could use log/slog instead here
	//lgr := slog.New(slog.NewTextHandler(lw, &slog.HandlerOptions{}))

	// ... lots of interesting things happen inside the application
	lgr.Print("Hello world")

	// Unix Logrotate can rename the file then send SIGHUP to the application.
	// If the application receives SIGHUP or SIGUSR1, "example.log" will be
	// closed then re-opened. So the old version is closed and kept and the
	// application carries on writing to "example.log", but it's now a new file,

	lgr.Print("Something interesting happened")
	// ... lots more interesting things happen inside the application

	// That's the end of the example. Now let's check that the log
	// file exists, with the correct messages in it.
	file, _ := os.Open("example.log")
	content, _ := io.ReadAll(file)
	fmt.Println(string(content))

	// Output: Hello world
	// Something interesting happened
}
