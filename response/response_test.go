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
	assert.Equal(t, "kong.response.get_status:null", getName(func() { response.GetStatus() }))
	res := make(chan int)
	go func(res chan int) { r, _ := response.GetStatus(); res <- r }(res)
	_ = <-ch
	ch <- "404"
	status := <-res
	assert.Equal(t, 404, status)
}

func TestGetHeaders(t *testing.T) {
	assert.Equal(t, "kong.response.get_headers:[100]", getName(func() { response.GetHeaders(100) }))
	assert.Equal(t, "kong.response.get_headers:null", getName(func() { response.GetHeaders(-1) }))

	res := make(chan map[string]interface{})
	go func(res chan map[string]interface{}) { r, _ := response.GetHeaders(-1); res <- r }(res)
	_ = <-ch
	ch <- `{"h1":"v1"}`
	headers := <-res
	assert.Equal(t, map[string]interface{}{"h1": "v1"}, headers)
}

func TestGetSource(t *testing.T) {
	assert.Equal(t, "kong.response.get_source:null", getName(func() { response.GetSource() }))
	assert.Equal(t, "foo", getStrValue(func(res chan string) { r, _ := response.GetSource(); res <- r }, "foo"))
	assert.Equal(t, "", getStrValue(func(res chan string) { r, _ := response.GetSource(); res <- r }, ""))
}

func TestSetStatus(t *testing.T) {
	assert.Equal(t, "kong.response.set_status:[404]", getName(func() { response.SetStatus(404) }))
}

func TestSetHeader(t *testing.T) {
	assert.Equal(t, `kong.response.set_header:["foo","bar"]`, getName(func() { response.SetHeader("foo", "bar") }))
}

func TestAddHeader(t *testing.T) {
	assert.Equal(t, `kong.response.add_header:["foo","bar"]`, getName(func() { response.AddHeader("foo", "bar") }))
}

func TestClearHeader(t *testing.T) {
	assert.Equal(t, `kong.response.clear_header:["foo"]`, getName(func() { response.ClearHeader("foo") }))
}

func TestSetHeaders(t *testing.T) {
	assert.Equal(t, `kong.response.set_headers:[{"h1":"v1"}]`, getName(func() {
		response.SetHeaders(map[string]interface{}{
			"h1": "v1",
		})
	}))
}
