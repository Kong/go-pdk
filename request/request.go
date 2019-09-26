package request

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	ch chan string
}

func NewRequest(ch chan string) *Request {
	return &Request{ch: ch}
}

func (r *Request) GetHeader(k string) string {
	r.ch <- fmt.Sprintf(`kong.request.get_header:%s`, k)
	return <-r.ch
}

func (r *Request) GetMethod() string {
	r.ch <- `kong.request.get_method`
	return <-r.ch
}

func (r *Request) GetQuery() map[string]interface{} {
	r.ch <- `kong.request.get_query`
	reply := <-r.ch
	query := make(map[string]interface{})
	json.Unmarshal([]byte(reply), &query)
	return query
}

func (r *Request) GetHeaders() map[string]interface{} {
	r.ch <- `kong.request.get_headers`
	reply := <-r.ch
	headers := make(map[string]interface{})
	json.Unmarshal([]byte(reply), &headers)
	return headers
}
