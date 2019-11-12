package ip

import (
	"github.com/Kong/go-pdk/bridge"
// 	"strconv"
)

type Ip struct {
	bridge.PdkBridge
}

func New(ch chan interface{}) Ip {
	return Ip{bridge.New(ch)}
}

func (ip Ip) IsTrusted(address string) (is_trusted bool, err error) {
	reply, err := ip.Ask(`kong.ip.is_trusted`, address)
	if err != nil {
		return
	}

	var ok bool
	if is_trusted, ok = reply.(bool); !ok {
		err = bridge.ReturnTypeError("boolean")
	}
	return
}
