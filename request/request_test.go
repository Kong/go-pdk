package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var request Request
var ch chan string

func init() {
	ch = make(chan string)
	request = New(ch)
}

func getName(f func()) string {
	go f()
	name := <-ch
	ch <- ""
	return name
}

func getStrValue(f func(res chan string), val string) string {
	res := make(chan string)
	go f(res)
	_ = <-ch
	ch <- val
	return <-res
}

func TestGetScheme(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetScheme() }), "kong.request.get_scheme:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetScheme(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetScheme(); res <- r }, ""), "")
}

func TestGetHost(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetHost() }), "kong.request.get_host:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetHost(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetHost(); res <- r }, ""), "")
}

func TestGetPort(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetPort() }), "kong.request.get_port:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetPort(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetPort(); res <- r }, ""), "")
}

func TestGetForwardedScheme(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetForwardedScheme() }), "kong.request.get_forwarded_scheme:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetForwardedScheme(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetForwardedScheme(); res <- r }, ""), "")
}

func TestGetForwardedHost(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetForwardedHost() }), "kong.request.get_forwarded_host:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetForwardedHost(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetForwardedHost(); res <- r }, ""), "")
}

func TestGetForwardedPort(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetForwardedPort() }), "kong.request.get_forwarded_port:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetForwardedPort(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetForwardedPort(); res <- r }, ""), "")
}

func TestGetHttpVersion(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetHttpVersion() }), "kong.request.get_http_version:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetHttpVersion(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetHttpVersion(); res <- r }, ""), "")
}

func TestGetMethod(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetMethod() }), "kong.request.get_method:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetMethod(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetMethod(); res <- r }, ""), "")
}

func TestGetPath(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetPath() }), "kong.request.get_path:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetPath(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetPath(); res <- r }, ""), "")
}

func TestGetPathWithQuery(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetPathWithQuery() }), "kong.request.get_path_with_query:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetPathWithQuery(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetPathWithQuery(); res <- r }, ""), "")
}

func TestGetRawQuery(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetRawQuery() }), "kong.request.get_raw_query:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetRawQuery(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetRawQuery(); res <- r }, ""), "")
}

func TestGetQueryArg(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetQueryArg() }), "kong.request.get_query_arg:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetQueryArg(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetQueryArg(); res <- r }, ""), "")
}

func TestGetQuery(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetQuery(100) }), "kong.request.get_query:[100]")
	assert.Equal(t, getName(func() { request.GetQuery(-1) }), "kong.request.get_query:null")

	res := make(chan map[string]interface{})
	go func(res chan map[string]interface{}) { r, _ := request.GetQuery(-1); res <- r }(res)
	_ = <-ch
	ch <- `{"h1":"v1"}`
	query := <-res
	assert.Equal(t, query, map[string]interface{}{"h1": "v1"})
}

func TestGetHeader(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetHeader("h1") }), `kong.request.get_header:["h1"]`)
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetHeader("h1"); res <- r }, "h1"), "h1")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := request.GetHeader("h2"); res <- r }, "h2"), "h2")
}

func TestGetHeaders(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetHeaders(100) }), "kong.request.get_headers:[100]")
	assert.Equal(t, getName(func() { request.GetHeaders(-1) }), "kong.request.get_headers:null")

	res := make(chan map[string]interface{})
	go func(res chan map[string]interface{}) { r, _ := request.GetHeaders(-1); res <- r }(res)
	_ = <-ch
	ch <- `{"h1":"v1"}`
	headers := <-res
	assert.Equal(t, headers, map[string]interface{}{"h1": "v1"})
}

func TestGetRawBody(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetRawBody() }), "kong.request.get_raw_body:null")

	res := make(chan string)
	go func(res chan string) { r, _ := request.GetRawBody(); res <- r }(res)
	_ = <-ch
	ch <- `{"h1":"v1"}`
	body := <-res
	assert.Equal(t, body, `{"h1":"v1"}`)
}
