package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var log Log
var ch chan string

func init() {
	ch = make(chan string)
	log = New(ch)
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

func TestAlert(t *testing.T) {
	assert.Equal(t, getName(func() { log.Alert("a", "b") }), "kong.log.alert:ab")
	assert.Equal(t, getName(func() { log.Alert("a", "b", "c") }), "kong.log.alert:abc")
}

func TestCrit(t *testing.T) {
	assert.Equal(t, getName(func() { log.Crit("a", "b") }), "kong.log.crit:ab")
	assert.Equal(t, getName(func() { log.Crit("a", "b", "c") }), "kong.log.crit:abc")
}

func TestErr(t *testing.T) {
	assert.Equal(t, getName(func() { log.Err("a", "b") }), "kong.log.err:ab")
	assert.Equal(t, getName(func() { log.Err("a", "b", "c") }), "kong.log.err:abc")
}

func TestWarn(t *testing.T) {
	assert.Equal(t, getName(func() { log.Warn("a", "b") }), "kong.log.warn:ab")
	assert.Equal(t, getName(func() { log.Warn("a", "b", "c") }), "kong.log.warn:abc")
}

func TestNotice(t *testing.T) {
	assert.Equal(t, getName(func() { log.Notice("a", "b") }), "kong.log.notice:ab")
	assert.Equal(t, getName(func() { log.Notice("a", "b", "c") }), "kong.log.notice:abc")
}

func TestInfo(t *testing.T) {
	assert.Equal(t, getName(func() { log.Info("a", "b") }), "kong.log.info:ab")
	assert.Equal(t, getName(func() { log.Info("a", "b", "c") }), "kong.log.info:abc")
}

func TestDebug(t *testing.T) {
	assert.Equal(t, getName(func() { log.Debug("a", "b") }), "kong.log.debug:ab")
	assert.Equal(t, getName(func() { log.Debug("a", "b", "c") }), "kong.log.debug:abc")
}

func TestSerialize(t *testing.T) {
	assert.Equal(t, getName(func() { log.Serialize() }), "kong.log.serialize")
	assert.Equal(t, getStrValue(func(res chan string) { res <- log.Serialize() }, ""), "")
	assert.Equal(t, getStrValue(func(res chan string) { res <- log.Serialize() }, "foo"), "foo")
}
