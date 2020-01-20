package service

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/stretchr/testify/assert"
)

var service Service
var ch chan interface{}

func init() {
	ch = make(chan interface{})
	service = New(ch)
}

func getBack(f func()) interface{} {
	go f()
	d := <-ch
	ch <- nil

	return d
}

func TestSetUpstream(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.service.set_upstream", Args: []interface{}{"foo"}}, getBack(func() { service.SetUpstream("foo") }))
}

func TestSetTarget(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.service.set_target", Args: []interface{}{"foo", 1}}, getBack(func() { service.SetTarget("foo", 1) }))
}
