package response

import (
// 	"strconv"

	"github.com/Kong/go-pdk/bridge"
)

type Response struct {
	bridge.PdkBridge
}

func New(ch chan interface{}) Response {
	return Response{bridge.New(ch)}
}

func (r Response) GetStatus() (i int, err error) {
	return r.AskInt(`kong.service.response.get_status`)
}

func (r Response) GetHeaders(max_headers int) (map[string]interface{}, error) {
	if max_headers == -1 {
		return r.AskMap(`kong.service.response.get_headers`)
	}

	return r.AskMap(`kong.service.response.get_headers`, max_headers)
}

func (r Response) GetHeader(name string) (string, error) {
	return r.AskString(`kong.service.response.get_header`, name)
}
