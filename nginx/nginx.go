package nginx

import (
	"fmt"
	"strconv"
	"github.com/kong/go-pdk/bridge"
)

type Nginx struct {
	bridge.PdkBridge
}

func New(ch chan string) Nginx {
	return Nginx{bridge.New(ch)}
}

func checkAndUnmarshall(s string, v interface{}) interface{} {
	if s == "null" {
		return nil
	}
	err := bridge.Unmarshal(s, &v)
	if err != nil {
		return nil
	}
	return v
}

func readStringP(reply string) *string {
	var s string
	if checkAndUnmarshall(reply, &s) == nil {
		return nil
	}
	return &s
}

func readFloatP(reply string) *float64 {
	var f float64
	if checkAndUnmarshall(reply, &f) == nil {
		return nil
	}
	return &f
}

func readAnyP(reply string) *interface{} {
	var p interface{}
	if checkAndUnmarshall(reply, &p) == nil {
		return nil
	}
	return &p
}

func (n Nginx) GetVar(k string) *string {
	return readStringP(n.Ask(fmt.Sprintf(`kong.nginx.get_var:%s`, k)))
}

func (n Nginx) GetTLS1VersionStr() *string {
	return readStringP(n.Ask(`kong.nginx.get_tls1_version_str`))
}

func (n Nginx) GetCtxAny(k string) *interface{} {
	return readAnyP(n.Ask(fmt.Sprintf(`kong.nginx.get_ctx:%s`, k)))
}

func (n Nginx) GetCtxString(k string) *string {
	return readStringP(n.Ask(fmt.Sprintf(`kong.nginx.get_ctx:%s`, k)))
}

func (n Nginx) GetCtxFloat(k string) *float64 {
	return readFloatP(n.Ask(fmt.Sprintf(`kong.nginx.get_ctx:%s`, k)))
}

func (n Nginx) ReqStartTime() float64 {
	reply := n.Ask(`kong.nginx.req_start_time`)
	t, _ := strconv.ParseFloat(reply, 64)
	return t
}
