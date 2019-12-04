package response

import (
	"testing"
	"github.com/Kong/go-pdk/bridge"
	"github.com/stretchr/testify/assert"
)

var response Response
var ch chan interface{}

func init() {
	ch = make(chan interface{})
	response = New(ch)
}

func getBack(f func()) interface{} {
	go f()
	d := <-ch
	ch <- nil

	return d
}

func TestGetStatus(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.service.response.get_status"}, getBack(func() { response.GetStatus() }))
}

func TestGetHeader(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.service.response.get_header", Args:[]interface{}{"foo"}}, getBack(func() { response.GetHeader("foo") }))
}

func TestGetHeaders(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.service.response.get_headers", Args:[]interface{}{1}}, getBack(func() { response.GetHeaders(1) }))
}
