package log

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/stretchr/testify/assert"
)

var log Log
var ch chan interface{}

func init() {
	ch = make(chan interface{})
	log = New(ch)
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

func TestAlert(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.log.alert", Args:[]interface{}{"Alo"}}, getBack(func() { log.Alert("Alo") }))
}

func TestCrit(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.log.crit", Args:[]interface{}{"Alo"}}, getBack(func() { log.Crit("Alo") }))
}

func TestErr(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.log.err", Args:[]interface{}{"Alo"}}, getBack(func() { log.Err("Alo") }))
}

func TestWarn(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.log.warn", Args:[]interface{}{"Alo"}}, getBack(func() { log.Warn("Alo") }))
}

func TestNotice(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.log.notice", Args:[]interface{}{"Alo"}}, getBack(func() { log.Notice("Alo") }))
}

func TestInfo(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.log.info", Args:[]interface{}{"Alo"}}, getBack(func() { log.Info("Alo") }))
}

func TestDebug(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.log.debug", Args:[]interface{}{"Alo"}}, getBack(func() { log.Debug("Alo") }))
}

func TestSerialize(t *testing.T) {
	assert.Equal(t, bridge.StepData{Method:"kong.log.serialize"}, getBack(func() { log.Serialize() }))
}
