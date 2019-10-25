package response

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	assert.Equal(t, getName(func() { response.GetStatus() }), "kong.response.get_status")
	res := make(chan int)
	go func(res chan int) { res <- response.GetStatus() }(res)
	_ = <-ch
	ch <- "404"
	status := <-res
	assert.Equal(t, status, 404)
}

func TestGetHeaders(t *testing.T) {
	assert.Equal(t, getName(func() { response.GetHeaders(100) }), "kong.response.get_headers:100")
	assert.Equal(t, getName(func() { response.GetHeaders(-1) }), "kong.response.get_headers")

	res := make(chan map[string]interface{})
	go func(res chan map[string]interface{}) { res <- response.GetHeaders(-1) }(res)
	_ = <-ch
	ch <- `{"h1":"v1"}`
	headers := <-res
	assert.Equal(t, headers, map[string]interface{}{"h1": "v1"})
}

func TestGetSource(t *testing.T) {
	assert.Equal(t, getName(func() { response.GetSource() }), "kong.response.get_source")
	assert.Equal(t, getStrValue(func(res chan string) { res <- response.GetSource() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- response.GetSource() }, ""), "")
}

func TestSetStatus(t *testing.T) {
	assert.Equal(t, getName(func() { response.SetStatus(404) }), "kong.response.set_status:404")
}

func TestSetHeader(t *testing.T) {
	assert.Equal(t, getName(func() { response.SetHeader("foo", "bar") }), `kong.response.set_header:["foo","bar"]`)
}

func TestAddHeader(t *testing.T) {
	assert.Equal(t, getName(func() { response.AddHeader("foo", "bar") }), `kong.response.add_header:["foo","bar"]`)
}

func TestClearHeader(t *testing.T) {
	assert.Equal(t, getName(func() { response.ClearHeader("foo") }), `kong.response.clear_header:foo`)
}

func TestSetHeaders(t *testing.T) {
	assert.Equal(t, getName(func() {
		response.SetHeaders(map[string]interface{}{
			"h1": "v1",
		})
	}), `kong.response.set_headers:{"h1":"v1"}`)
}
