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

func getStrValue(f func(res chan string), val string) string {
	res := make(chan string)
	go f(res)
	_ = <-ch
	ch <- val
	return <-res
}

func TestGetScheme(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_scheme"}, getBack(func() { request.GetScheme() }))
}

func TestGetHost(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_host"}, getBack(func() { request.GetHost() }))
}

func TestGetPort(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_port"}, getBack(func() { request.GetPort() }))
}

func TestGetForwardedScheme(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_scheme"}, getBack(func() { request.GetScheme() }))
}

func TestGetForwardedHost(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_forwarded_host"}, getBack(func() { request.GetForwardedHost() }))
}

func TestGetForwardedPort(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_forwarded_port"}, getBack(func() { request.GetForwardedPort() }))
}

func TestGetHttpVersion(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_http_version"}, getBack(func() { request.GetHttpVersion() }))
}

func TestGetMethod(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_method"}, getBack(func() { request.GetMethod() }))
}

func TestGetPath(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_path"}, getBack(func() { request.GetPath() }))
}

func TestGetPathWithQuery(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_path_with_query"}, getBack(func() { request.GetPathWithQuery() }))
}

func TestGetRawQuery(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_raw_query"}, getBack(func() { request.GetRawQuery() }))
}

func TestGetQueryArg(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_query_arg", Args:[]interface{}{"foo"}}, getBack(func() { request.GetQueryArg("foo") }))
}

func TestGetQuery(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_query", Args:[]interface{}{1}}, getBack(func() { request.GetQuery(1) }))
}

func TestGetHeader(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_header", Args:[]interface{}{"foo"}}, getBack(func() { request.GetHeader("foo") }))
}

func TestGetHeaders(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_headers", Args:[]interface{}{1}}, getBack(func() { request.GetHeaders(1) }))
}

func TestGetRawBody(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.request.get_raw_body"}, getBack(func() { request.GetRawBody() }))
}
