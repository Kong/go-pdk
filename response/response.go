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

func (r *Response) SetHeader(k string, v string) {
	r.ch <- fmt.Sprintf(`kong.response.set_header:["%s","%s"]`, k, v)
	_ = <-r.ch
}

func (r *Response) GetHeaders() map[string]interface{} {
	r.ch <- `kong.response.get_headers`
	reply := <-r.ch
	headers := make(map[string]interface{})
	json.Unmarshal([]byte(reply), &headers)
	return headers
}

func (r *Response) GetStatus() int {
	r.ch <- `kong.response.get_status`
	reply := <-r.ch
	status, _ := strconv.Atoi(reply)
	return status
}
