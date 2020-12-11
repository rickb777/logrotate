package safe

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestSafe(t *testing.T) {
	g := NewGomegaWithT(t)

	safe := New("foo")

	g.Expect(safe.Get()).Should(Equal("foo"))
	g.Expect(safe.Get()).Should(Not(Equal("bar")))

	safe.Put("bar")

	g.Expect(safe.Get()).Should(Equal("bar"))
	//g.Expect(safe.GetWhenDefined()).Should(Equal("bar"))

	safe.Put(nil)

	go func() {
		time.Sleep(time.Millisecond)
		safe.Put("yay")
	}()

	g.Expect(safe.GetWhenDefined()).Should(Equal("yay"))
}
