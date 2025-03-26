package safe

import (
	"github.com/rickb777/expect"
	"testing"
	"time"
)

func TestSafe(t *testing.T) {
	safe := New("foo")

	expect.Any(safe.Get()).ToBe(t, "foo")
	expect.Any(safe.Get()).Not().ToBe(t, "bar")

	safe.Put("bar")

	expect.Any(safe.Get()).ToBe(t, "bar")
	//expect.String(safe.GetWhenDefined()).ToBe(t,"bar"))

	safe.Put(nil)

	go func() {
		time.Sleep(time.Millisecond)
		safe.Put("yay")
	}()

	expect.Any(safe.GetWhenDefined()).ToBe(t, "yay")
}
