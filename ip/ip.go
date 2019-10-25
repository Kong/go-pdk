package ip

import (
	"fmt"
	"strconv"
	"github.com/kong/go-pdk/bridge"
)

type Ip struct {
	bridge.PdkBridge
}

func New(ch chan string) Ip {
	return Ip{bridge.New(ch)}
}

func (ip Ip) IsTrusted(address string) *bool {
	reply := ip.Ask(fmt.Sprintf(`kong.ip.is_trusted:%s`, address))
	is_trusted, err := strconv.ParseBool(reply)
	if err != nil {
		return nil
	}
	return &is_trusted
}
