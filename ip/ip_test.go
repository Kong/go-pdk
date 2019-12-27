package ip

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/stretchr/testify/assert"
)

var ip Ip
var ch chan interface{}

func init() {
	ch = make(chan interface{})
	ip = New(ch)
}

func getBack(f func()) interface{} {
	go f()
	d := <-ch
	ch <- nil

	return d
}

func TestIsTrusted(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.ip.is_trusted", Args: []interface{}{"1.1.1.1"}}, getBack(func() { ip.IsTrusted("1.1.1.1") }))
	assert.Equal(t, bridge.StepData{Method: "kong.ip.is_trusted", Args: []interface{}{"1.0.0.1"}}, getBack(func() { ip.IsTrusted("1.0.0.1") }))
}
