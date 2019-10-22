package response

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var response *Response
var ch chan string

func init() {
	ch = make(chan string)
	response = NewResponse(ch)
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

func TestGetStatus(t *testing.T) {
	assert.Equal(t, getName(func() { response.GetStatus() }), "kong.service.response.get_status")

	res := make(chan int)
	go func(res chan int) { res <- response.GetStatus() }(res)
	_ = <-ch
	ch <- "404"
	status := <-res
	assert.Equal(t, status, 404)
}

func TestGetHeader(t *testing.T) {
	assert.Equal(t, getName(func() { response.GetHeader("foo") }), "kong.service.response.get_header:foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- response.GetHeader("foo") }, "foo"), "foo")
}

func TestGetHeaders(t *testing.T) {
	assert.Equal(t, getName(func() { response.GetHeaders(100) }), "kong.service.response.get_headers:100")
	assert.Equal(t, getName(func() { response.GetHeaders(-1) }), "kong.service.response.get_headers")
}
