package router

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/stretchr/testify/assert"
)

var router Router
var ch chan interface{}

func init() {
	ch = make(chan interface{})
	router = New(ch)
}

func getBack(f func()) interface{} {
	go f()
	d := <-ch
	ch <- nil

	return d
}

func TestGetRoute(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.router.get_route"}, getBack(func() { router.GetRoute() }))
}

func TestGetService(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.router.get_service"}, getBack(func() { router.GetService() }))
}
