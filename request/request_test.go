package request

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var request *Request
var ch chan string

func init() {
	ch = make(chan string)
	request = &Request{ch: ch}
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
	assert.Equal(t, getName(func() { request.GetScheme() }), "kong.request.get_scheme")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetScheme() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetScheme() }, ""), "")
}

func TestGetHost(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetHost() }), "kong.request.get_host")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetHost() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetHost() }, ""), "")
}

func TestGetPort(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetPort() }), "kong.request.get_port")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetPort() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetPort() }, ""), "")
}

func TestGetForwardedScheme(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetForwardedScheme() }), "kong.request.get_forwarded_scheme")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetForwardedScheme() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetForwardedScheme() }, ""), "")
}

func TestGetForwardedHost(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetForwardedHost() }), "kong.request.get_forwarded_host")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetForwardedHost() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetForwardedHost() }, ""), "")
}

func TestGetForwardedPort(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetForwardedPort() }), "kong.request.get_forwarded_port")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetForwardedPort() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetForwardedPort() }, ""), "")
}

func TestGetHttpVersion(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetHttpVersion() }), "kong.request.get_http_version")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetHttpVersion() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetHttpVersion() }, ""), "")
}

func TestGetMethod(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetMethod() }), "kong.request.get_method")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetMethod() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetMethod() }, ""), "")
}

func TestGetPath(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetPath() }), "kong.request.get_path")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetPath() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetPath() }, ""), "")
}

func TestGetPathWithQuery(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetPathWithQuery() }), "kong.request.get_path_with_query")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetPathWithQuery() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetPathWithQuery() }, ""), "")
}

func TestGetRawQuery(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetRawQuery() }), "kong.request.get_raw_query")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetRawQuery() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetRawQuery() }, ""), "")
}

func TestGetQueryArg(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetQueryArg() }), "kong.request.get_query_arg")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetQueryArg() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetQueryArg() }, ""), "")
}

func TestGetQuery(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetQuery(100) }), "kong.request.get_query:100")
	assert.Equal(t, getName(func() { request.GetQuery(-1) }), "kong.request.get_query")

	res := make(chan map[string]interface{})
	go func(res chan map[string]interface{}) { res <- request.GetQuery(-1) }(res)
	_ = <-ch
	ch <- `{"h1":"v1"}`
	query := <-res
	assert.Equal(t, query, map[string]interface{}{"h1": "v1"})
}

func TestGetHeader(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetHeader("h1") }), "kong.request.get_header:h1")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetHeader("h1") }, "h1"), "h1")
	assert.Equal(t, getStrValue(func(res chan string) { res <- request.GetHeader("h2") }, "h2"), "h2")
}

func TestGetHeaders(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetHeaders(100) }), "kong.request.get_headers:100")
	assert.Equal(t, getName(func() { request.GetHeaders(-1) }), "kong.request.get_headers")

	res := make(chan map[string]interface{})
	go func(res chan map[string]interface{}) { res <- request.GetHeaders(-1) }(res)
	_ = <-ch
	ch <- `{"h1":"v1"}`
	headers := <-res
	assert.Equal(t, headers, map[string]interface{}{"h1": "v1"})
}

func TestGetRawBody(t *testing.T) {
	assert.Equal(t, getName(func() { request.GetRawBody() }), "kong.request.get_raw_body")

	res := make(chan string)
	go func(res chan string) { res <- request.GetRawBody() }(res)
	_ = <-ch
	ch <- `{"h1":"v1"}`
	body := <-res
	assert.Equal(t, body, `{"h1":"v1"}`)
}
