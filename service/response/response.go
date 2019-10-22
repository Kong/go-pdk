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
	r.ch <- `kong.service.response.get_status`
	reply := <-r.ch
	status, _ := strconv.Atoi(reply)
	return status
}

func (r *Response) GetHeaders(max_headers int) map[string]interface{} {
	var method string
	if max_headers == -1 {
		method = `kong.service.response.get_headers`
	} else {
		method = fmt.Sprintf(`kong.service.response.get_headers:%d`, max_headers)
	}
	r.ch <- method
	reply := <-r.ch
	headers := make(map[string]interface{})
	json.Unmarshal([]byte(reply), &headers)
	return headers
}

func (r *Response) GetHeader(name string) string {
	r.ch <- fmt.Sprintf(`kong.service.response.get_header:%s`, name)
	return <-r.ch
}
