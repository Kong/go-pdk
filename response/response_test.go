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
	assert.Equal(t, bridge.StepData{Method:"kong.response.get_status"}, getBack(func() { response.GetStatus() }))
}

func TestGetHeaders(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.response.get_headers", Args:[]interface{}{1}}, getBack(func() { response.GetHeaders(1) }))
}

func TestGetSource(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.response.get_source"}, getBack(func() { response.GetSource() }))
}

func TestSetStatus(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.response.get_status"}, getBack(func() { response.GetStatus() }))
}

func TestSetHeader(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.response.set_header", Args:[]interface{}{"foo", "bar"}}, getBack(func() { response.SetHeader("foo", "bar") }))
}

func TestAddHeader(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.response.add_header", Args:[]interface{}{"foo", "bar"}}, getBack(func() { response.AddHeader("foo", "bar") }))
}

func TestClearHeader(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.response.clear_header", Args:[]interface{}{"foo"}}, getBack(func() { response.ClearHeader("foo") }))
}

func TestSetHeaders(t *testing.T) {
	var m map[string]interface{} = nil
	assert.Equal(t, bridge.StepData{Method:"kong.response.set_headers", Args:[]interface{}{m}}, getBack(func() { response.SetHeaders(nil) }))
}
