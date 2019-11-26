package nginx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var nginx Nginx
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
	assert.Equal(t, `kong.nginx.get_var:["var1"]`, getName(func() { (nginx.GetVar("var1")) }))

	res := make(chan *string)
	go func(res chan *string) { r, _ := nginx.GetVar("var1"); res <- r }(res)
	_ = <-ch
	ch <- `"foo"`
	v := <-res
	assert.Equal(t, "foo", *v)
}

func TestGetTLS1VersionStr(t *testing.T) {
	assert.Equal(t, "kong.nginx.get_tls1_version_str:null", getName(func() { (nginx.GetTLS1VersionStr()) }))

	res := make(chan *string)
	go func(res chan *string) { r, _ := nginx.GetTLS1VersionStr(); res <- r }(res)
	_ = <-ch
	ch <- `"foo"`
	v := <-res
	assert.Equal(t, "foo", *v)
}

func TestGetCtxAny(t *testing.T) {
	assert.Equal(t, `kong.nginx.get_ctx:["v1"]`, getName(func() { (nginx.GetCtxAny("v1")) }))

	res := make(chan *interface{})
	go func(res chan *interface{}) { r, _ := nginx.GetCtxAny("v1"); res <- r }(res)
	_ = <-ch
	ch <- `"foo"`
	v := <-res
	assert.Equal(t, "foo", *v)
}

func TestGetCtxString(t *testing.T) {
	assert.Equal(t, `kong.nginx.get_ctx:["v1"]`, getName(func() { (nginx.GetCtxString("v1")) }))

	res := make(chan *string)
	go func(res chan *string) { r, _ := nginx.GetCtxString("v1"); res <- r }(res)
	_ = <-ch
	ch <- `"foo"`
	v := <-res
	assert.Equal(t, "foo", *v)
}

func TestGetCtxFloat(t *testing.T) {
	assert.Equal(t, `kong.nginx.get_ctx:["v1"]`, getName(func() { (nginx.GetCtxFloat("v1")) }))

	res := make(chan *float64)
	go func(res chan *float64) { r, _ := nginx.GetCtxFloat("v1"); res <- r }(res)
	_ = <-ch
	ch <- `12.4`
	v := <-res
	assert.Equal(t, 12.4, *v)
}

func TestReqStartTime(t *testing.T) {
	assert.Equal(t, "kong.nginx.req_start_time:null", getName(func() { (nginx.ReqStartTime()) }))

	res := make(chan float64)
	go func(res chan float64) { r, _ := nginx.ReqStartTime(); res <- r }(res)
	_ = <-ch
	ch <- `12.4`
	v := <-res
	assert.Equal(t, 12.4, v)
}
