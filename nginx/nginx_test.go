package nginx

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/stretchr/testify/assert"
)

var nginx Nginx
var ch chan interface{}

func init() {
	ch = make(chan interface{})
	nginx = New(ch)
}

func getName(f func()) interface{} {
	go f()
	name := <-ch
	ch <- ""
	return name
}

func getBack(f func()) interface{} {
	go f()
	d := <-ch
	ch <- nil

	return d
}

func TestGetVar(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.nginx.get_var", Args: []interface{}{"foo"}}, getBack(func() { nginx.GetVar("foo") }))
}

func TestGetTLS1VersionStr(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.nginx.get_tls1_version_str"}, getBack(func() { nginx.GetTLS1VersionStr() }))
}

func TestGetCtxAny(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.nginx.get_ctx", Args: []interface{}{"foo"}}, getBack(func() { nginx.GetCtxAny("foo") }))
}

func TestGetCtxString(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.nginx.get_ctx", Args: []interface{}{"foo"}}, getBack(func() { nginx.GetCtxString("foo") }))
}

func TestGetCtxFloat(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.nginx.get_ctx", Args: []interface{}{"foo"}}, getBack(func() { nginx.GetCtxFloat("foo") }))
}

func TestReqStartTime(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method: "kong.nginx.get_ctx", Args: []interface{}{"foo"}}, getBack(func() { nginx.GetCtxFloat("foo") }))
}
