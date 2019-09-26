package nginx

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Nginx struct {
	ch chan string
}

func NewNginx(ch chan string) *Nginx {
	return &Nginx{ch: ch}
}

func readStringP(ch chan string) *string {
	reply := <-ch
	if reply == "null" {
		return nil
	}
	var s string
	err := json.Unmarshal([]byte(reply), &s)
	if err != nil {
		return nil
	}
	return &s
}

func readFloatP(ch chan string) *float64 {
	reply := <-ch
	if reply == "null" {
		return nil
	}
	var f float64
	err := json.Unmarshal([]byte(reply), &f)
	if err != nil {
		return nil
	}
	return &f
}

func readAnyP(ch chan string) *interface{} {
	reply := <-ch
	if reply == "null" {
		return nil
	}
	var p interface{}
	err := json.Unmarshal([]byte(reply), &p)
	if err != nil {
		return nil
	}
	return &p
}

func (n *Nginx) GetVar(k string) *string {
	n.ch <- fmt.Sprintf(`kong.nginx.get_var:%s`, k)
	return readStringP(n.ch)
}

func (n *Nginx) GetTLS1VersionStr() *string {
	n.ch <- `kong.nginx.get_tls1_version_str:%s`
	return readStringP(n.ch)
}

func (n *Nginx) GetCtxAny(k string) *interface{} {
	n.ch <- fmt.Sprintf(`kong.nginx.get_ctx:%s`, k)
	return readAnyP(n.ch)
}

func (n *Nginx) GetCtxString(k string) *string {
	n.ch <- fmt.Sprintf(`kong.nginx.get_ctx:%s`, k)
	return readStringP(n.ch)
}

func (n *Nginx) GetCtxFloat(k string) *float64 {
	n.ch <- fmt.Sprintf(`kong.nginx.get_ctx:%s`, k)
	return readFloatP(n.ch)
}

func (n *Nginx) ReqStartTime() float64 {
	n.ch <- `kong.nginx.req_start_time`
	reply := <-n.ch
	t, _ := strconv.ParseFloat(reply, 64)
	return t
}
