package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var service *Service
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
	assert.Equal(t, getName(func() { service.SetUpstream("example.test") }), "kong.service.set_upstream:example.test")
}

func TestSetTarget(t *testing.T) {
	assert.Equal(t, getName(func() { service.SetTarget("example.test", 80) }), `kong.service.set_target:["example.test", 80]`)
}
