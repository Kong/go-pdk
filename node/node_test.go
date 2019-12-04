package node

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/stretchr/testify/assert"
)

var node Node
var ch chan interface{}

func init() {
	ch = make(chan interface{})
	node = New(ch)
}

func getBack(f func()) interface{} {
	go f()
	d := <-ch
	ch <- nil

	return d
}

func TestGetId(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.node.get_id"}, getBack(func() { node.GetId() }))
}

func TestGetMemoryStats(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.node.get_memory_stats"}, getBack(func() { node.GetMemoryStats() }))
}
