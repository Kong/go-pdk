package request

import (
	"fmt"
	"github.com/kong/go-pdk/bridge"
)

type Request struct {
	bridge.PdkBridge
}

func New(ch chan string) Request {
	return Request{bridge.New(ch)}
}

func (r Request) SetScheme(scheme string) {
	_ = r.Ask(fmt.Sprintf(`kong.service.request.set_scheme:%s`, scheme))
}

func (r Request) SetPath(path string) {
	_ = r.Ask(fmt.Sprintf(`kong.service.request.set_path:%s`, path))
}

func (r Request) SetRawQuery(query string) {
	_ = r.Ask(fmt.Sprintf(`kong.service.request.set_raw_query:%s`, query))
}

func (r Request) SetMethod(method string) {
	_ = r.Ask(fmt.Sprintf(`kong.service.request.set_method:%s`, method))
}

func (r Request) SetQuery(query string) {
	_ = r.Ask(fmt.Sprintf(`kong.service.request.set_query:%s`, query))
}

func (r Request) SetHeader(name string, value string) {
	_ = r.Ask(fmt.Sprintf(`kong.service.request.set_header:["%s", "%s"]`, name, value))
}

func (r Request) AddHeader(name string, value string) {
	_ = r.Ask(fmt.Sprintf(`kong.service.request.add_header:["%s", "%s"]`, name, value))
}

func (r Request) ClearHeader(name string) {
	_ = r.Ask(fmt.Sprintf(`kong.service.request.clear_header:%s`, name))
}

func (r Request) SetHeaders(headers map[string]interface{}) error {
	headersBytes, err := bridge.Marshal(headers)
	if err != nil {
		return err
	}

	_ = r.Ask(fmt.Sprintf(`kong.service.request.set_headers:%s`, headersBytes))
	return nil
}

func (r Request) SetRawBody(body string) {
	_ = r.Ask(fmt.Sprintf(`kong.service.request.set_raw_body:%s`, body))
}

// TODO set_body
