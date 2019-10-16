package ip

import (
	"strconv"
)

type Ip struct {
	ch chan string
}

func NewIp(ch chan string) *Ip {
	return &Ip{ch: ch}
}

func (ip *Ip) IsTrusted(address string) *bool {
	ip.ch <- `kong.ip.is_trusted`
	reply := <- ip.ch
	is_trusted, err := strconv.ParseBool(reply)
	if err != nil {
		return nil
	}
	return &is_trusted
}
