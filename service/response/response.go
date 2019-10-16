package response

import (
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

// TODO get_headers

func (r *Response) GetHeader(name string) string {
	r.ch <- fmt.Sprintf(`kong.service.response.get_header:%s`, name)
	return <-r.ch
}

