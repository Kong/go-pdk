package response

import (
	"fmt"
	"strconv"
	"github.com/kong/go-pdk/bridge"
)

type Response struct {
	bridge.PdkBridge
}

func New(ch chan string) *Response {
	return &Response{*bridge.New(ch)}
}

func (r *Response) GetStatus() int {
	reply := r.Ask(`kong.response.get_status`)
	status, _ := strconv.Atoi(reply)
	return status
}

func (r *Response) GetHeader(name string) string {
	return r.Ask(fmt.Sprintf(`kong.response.get_header:%s`, name))
}

func (r *Response) GetHeaders(max_headers int) map[string]interface{} {
	var method string
	if max_headers == -1 {
		method = `kong.response.get_headers`
	} else {
		method = fmt.Sprintf(`kong.response.get_headers:%d`, max_headers)
	}

	headers := make(map[string]interface{})
	bridge.Unmarshal(r.Ask(method), &headers)
	return headers
}

func (r *Response) GetSource() string {
	return r.Ask(`kong.response.get_source`)
}

func (r *Response) SetStatus(status int) {
	_ = r.Ask(fmt.Sprintf(`kong.response.set_status:%d`, status))
}

func (r *Response) SetHeader(k string, v string) {
	_ = r.Ask(fmt.Sprintf(`kong.response.set_header:["%s","%s"]`, k, v))
}

func (r *Response) AddHeader(k string, v string) {
	_ = r.Ask(fmt.Sprintf(`kong.response.add_header:["%s","%s"]`, k, v))
}

func (r *Response) ClearHeader(k string) {
	_ = r.Ask(fmt.Sprintf(`kong.response.clear_header:%s`, k))
}

func (r *Response) SetHeaders(headers map[string]interface{}) error {
	headersBytes, err := bridge.Marshal(headers)
	if err != nil {
		return err
	}

	_ = r.Ask(fmt.Sprintf(`kong.response.set_headers:%s`, headersBytes))
	return nil
}

// TODO exit
