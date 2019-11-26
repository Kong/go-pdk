package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var service Service
var ch chan string

func init() {
	ch = make(chan string)
	service = New(ch)
}

func getName(f func()) string {
	go f()
	name := <-ch
	ch <- ""
	return name
}

func TestSetUpstream(t *testing.T) {
	assert.Equal(t, `kong.service.set_upstream:["example.test"]`, getName(func() { service.SetUpstream("example.test") }))
}

func TestSetTarget(t *testing.T) {
	assert.Equal(t, `kong.service.set_target:["example.test",80]`, getName(func() { service.SetTarget("example.test", 80) }))
}
