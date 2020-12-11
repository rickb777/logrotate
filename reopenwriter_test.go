package logrotate

import (
	. "github.com/onsi/gomega"
	ioioutil "io/ioutil"
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

	content, err := ioioutil.ReadAll(file)
	g.Expect(err).Should(BeNil())

	g.Expect(file.Close()).Should(BeNil())

	g.Expect(string(content)).Should(Equal(line1 + line2))
}
