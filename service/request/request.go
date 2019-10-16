package request

import (
	"fmt"
)

type Request struct {
	ch chan string
}

func NewRequest(ch chan string) *Request {
	return &Request{ch: ch}
}

func (r *Request) SetScheme(scheme string) {
	r.ch <- fmt.Sprintf(`kong.request.response.set_scheme:%s`, scheme)
	_ = <- r.ch
}

func (r *Request) SetPath(path string) {
	r.ch <- fmt.Sprintf(`kong.request.response.set_path:%s`, path)
	_ = <- r.ch
}

func (r *Request) SetRawQuery(query string) {
	r.ch <- fmt.Sprintf(`kong.request.response.set_raw_query:%s`, query)
	_ = <- r.ch
}

func (r *Request) SetMethod(method string) {
	r.ch <- fmt.Sprintf(`kong.request.response.set_method:%s`, method)
	_ = <- r.ch
}

func (r *Request) SetQuery(query string) {
	r.ch <- fmt.Sprintf(`kong.request.response.set_query:%s`, query)
	_ = <- r.ch
}

func (r *Request) Setheader(name string, value string) {
	r.ch <- fmt.Sprintf(`kong.request.response.set_header:["%s", "%s"]`, name, value)
	_ = <- r.ch
}

func (r *Request) AddHeader(name string, value string) {
	r.ch <- fmt.Sprintf(`kong.request.response.add_header:["%s", "%s"]`, name, value)
	_ = <- r.ch
}

func (r *Request) ClearHeader(name string) {
	r.ch <- fmt.Sprintf(`kong.request.response.clear_header:%s`, name)
	_ = <- r.ch
}

// TODO set_headers

func (r *Request) SetRawBody(body string) {
	r.ch <- fmt.Sprintf(`kong.request.response.set_raw_body:%s`, body)
	_ = <- r.ch
}

// TODO set_body
