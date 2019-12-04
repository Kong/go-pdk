package client

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/entities"
)

var client Client
var ch chan interface {}

func init() {
	ch = make(chan interface {})
	client = New(ch)
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

func TestGetIp(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.client.get_ip"}, getBack(func() { client.GetIp() }))
	assert.Equal(t, "foo", getStrValue(func(res chan string) { r, _ := client.GetIp(); res <- r }, "foo"))
	assert.Equal(t, "", getStrValue(func(res chan string) { r, _ := client.GetIp(); res <- r }, ""))
}

func TestGetForwardedIp(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.client.get_forwarded_ip"}, getBack(func() { client.GetForwardedIp() }))
	assert.Equal(t, "foo", getStrValue(func(res chan string) { r, _ := client.GetForwardedIp(); res <- r }, "foo"))
	assert.Equal(t, "", getStrValue(func(res chan string) { r, _ := client.GetForwardedIp(); res <- r }, ""))
}

func TestGetPort(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.client.get_port"}, getBack(func() { client.GetPort() }))
	assert.Equal(t, "foo", getStrValue(func(res chan string) { r, _ := client.GetPort(); res <- r }, "foo"))
	assert.Equal(t, "", getStrValue(func(res chan string) { r, _ := client.GetPort(); res <- r }, ""))
}

func TestGetForwardedPort(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.client.get_forwarded_port"}, getBack(func() { client.GetForwardedPort() }))
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := client.GetForwardedPort(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := client.GetForwardedPort(); res <- r }, ""), "")
}

func TestGetCredential(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.client.get_credential"}, getBack(func() { client.GetCredential() }))
}

func TestLoadConsumer(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.client.load_consumer", Args:[]interface{}{"foo", true}}, getBack(func() { client.LoadConsumer("foo", true) }))
}

func TestGetConsumer(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.client.get_consumer"}, getBack(func() { client.GetConsumer() }))
}

func TestAuthenticate(t *testing.T) {
	var consumer *entities.Consumer = nil
	var credential *AuthenticatedCredential = nil
	assert.Equal(t, bridge.StepData{Method:"kong.client.authenticate", Args:[]interface{}{consumer, credential}}, getBack(func() { client.Authenticate(nil, nil) }))
}

func TestGetProtocol(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.client.get_protocol", Args:[]interface{}{true}}, getBack(func() { client.GetProtocol(true) }))
	assert.Equal(t, bridge.StepData{Method:"kong.client.get_protocol", Args:[]interface{}{false}}, getBack(func() { client.GetProtocol(false) }))
}
