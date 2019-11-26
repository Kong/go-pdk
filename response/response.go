package response

import (
	"strconv"

	"github.com/Kong/go-pdk/bridge"
)

type Response struct {
	bridge.PdkBridge
}

func New(ch chan string) Response {
	return Response{bridge.New(ch)}
}

func (r Response) GetStatus() (int, error) {
	reply, err := r.Ask(`kong.response.get_status`)
	if err != nil {
		return 0, err
	}

	status, _ := strconv.Atoi(reply)
	return status, nil
}

func (r Response) GetHeader(name string) (string, error) {
	return r.Ask(`kong.response.get_header`, name)
}

func (r Response) GetHeaders(max_headers int) (map[string]interface{}, error) {
	var res string
	var err error
	if max_headers == -1 {
		res, err = r.Ask(`kong.response.get_headers`)
	} else {
		res, err = r.Ask(`kong.response.get_headers`, max_headers)
	}
	if err != nil {
		return nil, err
	}

	headers := make(map[string]interface{})
	bridge.Unmarshal(res, &headers)
	return headers, nil
}

func (r Response) GetSource() (string, error) {
	return r.Ask(`kong.response.get_source`)
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
