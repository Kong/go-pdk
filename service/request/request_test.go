package request

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/stretchr/testify/assert"
)

var request Request
var ch chan interface{}

func init() {
	ch = make(chan interface{})
	request = New(ch)
}

func getBack(f func()) interface{} {
	go f()
	d := <-ch
	ch <- nil

	return d
}

func TestSetScheme(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.service.request.set_scheme", Args: []interface{}{"http"}}, getBack(func() { request.SetScheme("http") }))
}

func TestSetPath(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.service.request.set_path", Args: []interface{}{"/foo"}}, getBack(func() { request.SetPath("/foo") }))
}

func TestSetRawQuery(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.service.request.set_raw_query", Args: []interface{}{"name=foo"}}, getBack(func() { request.SetRawQuery("name=foo") }))
}

func TestSetMethod(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.service.request.set_method", Args: []interface{}{"GET"}}, getBack(func() { request.SetMethod("GET") }))
}

func TestSetQuery(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.service.request.set_query", Args: []interface{}{map[string][]string{"foo": {"bar"}}}}, getBack(func() { request.SetQuery(map[string][]string{"foo": {"bar"}}) }))
}

func TestSetHeader(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.service.request.set_header", Args: []interface{}{"foo", "bar"}}, getBack(func() { request.SetHeader("foo", "bar") }))
}

func TestAddHeader(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.service.request.add_header", Args: []interface{}{"foo", "bar"}}, getBack(func() { request.AddHeader("foo", "bar") }))
}

func TestClearHeader(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.service.request.clear_header", Args: []interface{}{"foo"}}, getBack(func() { request.ClearHeader("foo") }))
}

func TestSetHeaders(t *testing.T) {
	var h map[string][]string = nil
	assert.Equal(t, bridge.StepData{Method: "kong.service.request.set_headers", Args: []interface{}{h}}, getBack(func() { request.SetHeaders(nil) }))
}

func TestSetRawBody(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.service.request.set_raw_body", Args: []interface{}{"foo"}}, getBack(func() { request.SetRawBody("foo") }))
}
