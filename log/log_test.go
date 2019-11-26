package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, `kong.log.alert:["a","b"]`, getName(func() { log.Alert("a", "b") }))
	assert.Equal(t, `kong.log.alert:["a","b","c"]`, getName(func() { log.Alert("a", "b", "c") }))
}

func TestCrit(t *testing.T) {
	assert.Equal(t, `kong.log.crit:["a","b"]`, getName(func() { log.Crit("a", "b") }))
	assert.Equal(t, `kong.log.crit:["a","b","c"]`, getName(func() { log.Crit("a", "b", "c") }))
}

func TestErr(t *testing.T) {
	assert.Equal(t, `kong.log.err:["a","b"]`, getName(func() { log.Err("a", "b") }))
	assert.Equal(t, `kong.log.err:["a","b","c"]`, getName(func() { log.Err("a", "b", "c") }))
}

func TestWarn(t *testing.T) {
	assert.Equal(t, `kong.log.warn:["a","b"]`, getName(func() { log.Warn("a", "b") }))
	assert.Equal(t, `kong.log.warn:["a","b","c"]`, getName(func() { log.Warn("a", "b", "c") }))
}

func TestNotice(t *testing.T) {
	assert.Equal(t, `kong.log.notice:["a","b"]`, getName(func() { log.Notice("a", "b") }))
	assert.Equal(t, `kong.log.notice:["a","b","c"]`, getName(func() { log.Notice("a", "b", "c") }))
}

func TestInfo(t *testing.T) {
	assert.Equal(t, `kong.log.info:["a","b"]`, getName(func() { log.Info("a", "b") }))
	assert.Equal(t, `kong.log.info:["a","b","c"]`, getName(func() { log.Info("a", "b", "c") }))
}

func TestDebug(t *testing.T) {
	assert.Equal(t, `kong.log.debug:["a","b"]`, getName(func() { log.Debug("a", "b") }))
	assert.Equal(t, `kong.log.debug:["a","b","c"]`, getName(func() { log.Debug("a", "b", "c") }))
}

func TestSerialize(t *testing.T) {
	assert.Equal(t, "kong.log.serialize:null", getName(func() { log.Serialize() }))
	assert.Equal(t, "", getStrValue(func(res chan string) { r, _ := log.Serialize(); res <- r }, ""))
	assert.Equal(t, "foo", getStrValue(func(res chan string) { r, _ := log.Serialize(); res <- r }, "foo"))
}
