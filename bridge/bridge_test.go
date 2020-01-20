package bridge

import (
	"testing"

	"github.com/Kong/go-pdk/entities"
	"github.com/stretchr/testify/assert"
)

var ch chan interface{}
var bridge PdkBridge

func init() {
	ch = make(chan interface{})
	bridge = New(ch)
}

func TestAsk(t *testing.T) {
	go func() {
		bridge.Ask("foo.bar", 1, 2, 3, 1.23, false)
	}()

	call := <-ch
	ch <- ""

	assert.Equal(t, call, StepData{
		Method: "foo.bar",
		Args:   []interface{}{1, 2, 3, 1.23, false},
	})

	go func() {
		n := "gs"
		bridge.Ask("foo.bar", entities.Consumer{Username: n})
	}()

	call = <-ch
	ch <- ""

	n := "gs"
	consumer := []interface{}{entities.Consumer{Username: n}}
	assert.Equal(t, StepData{
		Method: "foo.bar",
		Args:   consumer,
	}, call)
}
