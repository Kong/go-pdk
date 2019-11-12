package response

import (
	"github.com/Kong/go-pdk/bridge"
)

type Response struct {
	bridge.PdkBridge
}

func New(ch chan interface{}) Response {
	return Response{bridge.New(ch)}
}

func (r Response) GetStatus() (int, error) {
	return r.AskInt(`kong.response.get_status`)
}

func (r Response) GetHeader(name string) (string, error) {
	return r.AskString(`kong.response.get_header`, name)
}

func (r Response) GetHeaders(max_headers int) (res map[string]interface{}, err error) {
	if max_headers == -1 {
		return r.AskMap(`kong.response.get_headers`)
	}

	return r.AskMap(`kong.response.get_headers`, max_headers)
}

func (r Response) GetSource() (string, error) {
	return r.AskString(`kong.response.get_source`)
}

func (r Response) SetStatus(status int) error {
	_, err := r.Ask(`kong.response.set_status`, status)
	return err
}

func (r Response) SetHeader(k string, v string) error {
	_, err := r.Ask(`kong.response.set_header`, k, v)
	return err
}

func (r Response) AddHeader(k string, v string) error {
	_, err := r.Ask(`kong.response.add_header`, k, v)
	return err
}

func (r Response) ClearHeader(k string) error {
	_, err := r.Ask(`kong.response.clear_header`, k)
	return err
}

func (r Response) SetHeaders(headers map[string]interface{}) error {
	_, err := r.Ask(`kong.response.set_headers`, headers)
	return err
}

// TODO exit
