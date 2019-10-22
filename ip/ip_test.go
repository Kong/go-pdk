package ip

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var ip *Ip
var ch chan string

func init() {
	ch = make(chan string)
	ip = &Ip{ch: ch}
}

func getName(f func()) string {
	go f()
	name := <-ch
	ch <- ""
	return name
}

func TestIsTrusted(t *testing.T) {
	assert.Equal(t, getName(func() { ip.IsTrusted("1.1.1.1") }), "kong.ip.is_trusted:1.1.1.1")
	assert.Equal(t, getName(func() { ip.IsTrusted("1.0.0.1") }), "kong.ip.is_trusted:1.0.0.1")

	res := make(chan *bool)
	go func(res chan *bool) { res <- ip.IsTrusted("1.1.1.1") }(res)
	_ = <-ch
	ch <- `true`
	trusted := <-res
	assert.Equal(t, *trusted, true)
}
