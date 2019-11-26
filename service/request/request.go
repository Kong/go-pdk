package request

import (
	"github.com/Kong/go-pdk/bridge"
)

type Request struct {
	bridge.PdkBridge
}

func New(ch chan string) Request {
	return Request{bridge.New(ch)}
}

func (r Request) SetScheme(scheme string) error {
	_, err := r.Ask(`kong.service.request.set_scheme`, scheme)
	return err
}

func (r Request) SetPath(path string) error {
	_, err := r.Ask(`kong.service.request.set_path`, path)
	return err
}

func (r Request) SetRawQuery(query string) error {
	_, err := r.Ask(`kong.service.request.set_raw_query`, query)
	return err
}

func (r Request) SetMethod(method string) error {
	_, err := r.Ask(`kong.service.request.set_method`, method)
	return err
}

func (r Request) SetQuery(query string) error {
	_, err := r.Ask(`kong.service.request.set_query`, query)
	return err
}

func (r Request) SetHeader(name string, value string) error {
	_, err := r.Ask(`kong.service.request.set_header`, name, value)
	return err
}

func (r Request) AddHeader(name string, value string) error {
	_, err := r.Ask(`kong.service.request.add_header`, name, value)
	return err
}

func (r Request) ClearHeader(name string) error {
	_, err := r.Ask(`kong.service.request.clear_header`, name)
	return err
}

func (r Request) SetHeaders(headers map[string]interface{}) error {
	_, err := r.Ask(`kong.service.request.set_headers`, headers)
	return err
}

func (r Request) SetRawBody(body string) error {
	_, err := r.Ask(`kong.service.request.set_raw_body`, body)
	return err
}

// TODO set_body
