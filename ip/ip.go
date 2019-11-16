package ip

import (
	"github.com/Kong/go-pdk/bridge"
	"strconv"
)

type Ip struct {
	bridge.PdkBridge
}

func New(ch chan string) Ip {
	return Ip{bridge.New(ch)}
}

func (ip Ip) IsTrusted(address string) (*bool, error) {
	reply, err := ip.Ask(`kong.ip.is_trusted`, address)
	if err != nil {
		return nil, err
	}

	is_trusted, err := strconv.ParseBool(reply)
	if err != nil {
		return nil, err
	}
	return &is_trusted, nil
}
