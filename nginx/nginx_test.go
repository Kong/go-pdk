package nginx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var nginx *Nginx
var ch chan string

func init() {
	ch = make(chan string)
	nginx = New(ch)
}

func getName(f func()) string {
	go f()
	name := <-ch
	ch <- ""
	return name
}

func TestGetVar(t *testing.T) {
	assert.Equal(t, getName(func() { (nginx.GetVar("var1")) }), "kong.nginx.get_var:var1")

	res := make(chan *string)
	go func(res chan *string) { res <- nginx.GetVar("var1") }(res)
	_ = <-ch
	ch <- `"foo"`
	v := <-res
	assert.Equal(t, *v, "foo")
}

func TestGetTLS1VersionStr(t *testing.T) {
	assert.Equal(t, getName(func() { (nginx.GetTLS1VersionStr()) }), "kong.nginx.get_tls1_version_str")

	res := make(chan *string)
	go func(res chan *string) { res <- nginx.GetTLS1VersionStr() }(res)
	_ = <-ch
	ch <- `"foo"`
	v := <-res
	assert.Equal(t, *v, "foo")
}

func TestGetCtxAny(t *testing.T) {
	assert.Equal(t, getName(func() { (nginx.GetCtxAny("v1")) }), "kong.nginx.get_ctx:v1")

	res := make(chan *interface{})
	go func(res chan *interface{}) { res <- nginx.GetCtxAny("v1") }(res)
	_ = <-ch
	ch <- `"foo"`
	v := <-res
	assert.Equal(t, *v, "foo")
}

func TestGetCtxString(t *testing.T) {
	assert.Equal(t, getName(func() { (nginx.GetCtxString("v1")) }), "kong.nginx.get_ctx:v1")

	res := make(chan *string)
	go func(res chan *string) { res <- nginx.GetCtxString("v1") }(res)
	_ = <-ch
	ch <- `"foo"`
	v := <-res
	assert.Equal(t, *v, "foo")
}

func TestGetCtxFloat(t *testing.T) {
	assert.Equal(t, getName(func() { (nginx.GetCtxFloat("v1")) }), "kong.nginx.get_ctx:v1")

	res := make(chan *float64)
	go func(res chan *float64) { res <- nginx.GetCtxFloat("v1") }(res)
	_ = <-ch
	ch <- `12.4`
	v := <-res
	assert.Equal(t, *v, 12.4)
}

func TestReqStartTime(t *testing.T) {
	assert.Equal(t, getName(func() { (nginx.ReqStartTime()) }), "kong.nginx.req_start_time")

	res := make(chan float64)
	go func(res chan float64) { res <- nginx.ReqStartTime() }(res)
	_ = <-ch
	ch <- `12.4`
	v := <-res
	assert.Equal(t, v, 12.4)
}
