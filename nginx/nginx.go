package nginx

import (
	"strconv"

	"github.com/Kong/go-pdk/bridge"
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

func (n Nginx) GetVar(k string) (*string, error) {
	if res, err := n.Ask(`kong.nginx.get_var`, k); err != nil {
		return nil, err
	} else {
		return readStringP(res), nil
	}
}

func (n Nginx) GetTLS1VersionStr() (*string, error) {
	if res, err := n.Ask(`kong.nginx.get_tls1_version_str`); err != nil {
		return nil, err
	} else {
		return readStringP(res), nil
	}
}

func (n Nginx) GetCtxAny(k string) (*interface{}, error) {
	if res, err := n.Ask(`kong.nginx.get_ctx`, k); err != nil {
		return nil, err
	} else {
		return readAnyP(res), nil
	}
}

func (n Nginx) GetCtxString(k string) (*string, error) {
	if res, err := n.Ask(`kong.nginx.get_ctx`, k); err != nil {
		return nil, err
	} else {
		return readStringP(res), nil
	}
}

func (n Nginx) GetCtxFloat(k string) (*float64, error) {
	if res, err := n.Ask(`kong.nginx.get_ctx`, k); err != nil {
		return nil, err
	} else {
		return readFloatP(res), nil
	}
}

func (n Nginx) ReqStartTime() (float64, error) {
	if res, err := n.Ask(`kong.nginx.req_start_time`); err != nil {
		return 0, err
	} else {
		t, _ := strconv.ParseFloat(res, 64)
		return t, nil
	}
}
