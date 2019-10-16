package response

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Response struct {
	ch chan string
}

func NewResponse(ch chan string) *Response {
	return &Response{ch: ch}
}

func (r *Response) GetStatus() int {
	r.ch <- `kong.response.get_status`
	reply := <-r.ch
	status, _ := strconv.Atoi(reply)
	return status
}

func (r *Response) GetHeader(name string) string {
	r.ch <- fmt.Sprintf(`kong.response.get_header:%s`, name)
	return <-r.ch
}

func (r *Response) GetHeaders() map[string]interface{} {
	r.ch <- `kong.response.get_headers`
	reply := <-r.ch
	headers := make(map[string]interface{})
	json.Unmarshal([]byte(reply), &headers)
	return headers
}

func (r *Response) GetSource() string {
	r.ch <- `kong.response.get_source`
	return <-r.ch
}

func (r *Response) SetStatus(status int) {
	r.ch <- fmt.Sprintf(`kong.response.set_status:%d`, status)
	_ = <-r.ch
}

func (r *Response) SetHeader(k string, v string) {
	r.ch <- fmt.Sprintf(`kong.response.set_header:["%s","%s"]`, k, v)
	_ = <-r.ch
}

func (r *Response) AddHeader(k string, v string) {
	r.ch <- fmt.Sprintf(`kong.response.add_header:["%s","%s"]`, k, v)
	_ = <-r.ch
}

func (r *Response) ClearHeader(k string) {
	r.ch <- fmt.Sprintf(`kong.response.clear_header:%s`, k)
	_ = <-r.ch
}

// TODO set_headers

// TODO exit
