package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var response Response
var ch chan string

func init() {
	ch = make(chan string)
	response = New(ch)
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
	assert.Equal(t, "kong.service.response.get_status:null", getName(func() { response.GetStatus() }))

	res := make(chan int)
	go func(res chan int) { r, _ := response.GetStatus(); res <- r }(res)
	_ = <-ch
	ch <- "404"
	status := <-res
	assert.Equal(t, 404, status)
}

func TestGetHeader(t *testing.T) {
	assert.Equal(t, `kong.service.response.get_header:["foo"]`, getName(func() { response.GetHeader("foo") }))
	assert.Equal(t, "foo", getStrValue(func(res chan string) { r, _ := response.GetHeader("foo"); res <- r }, "foo"))
}

func TestGetHeaders(t *testing.T) {
	assert.Equal(t, "kong.service.response.get_headers:[100]", getName(func() { response.GetHeaders(100) }))
	assert.Equal(t, "kong.service.response.get_headers:null", getName(func() { response.GetHeaders(-1) }))
}
