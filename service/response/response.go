package response

import (
	"fmt"
	"strconv"
	"github.com/kong/go-pdk/bridge"
)

type Response struct {
	bridge.PdkBridge
}

func New(ch chan string) Response{
	return Response{bridge.New(ch)}
}

func (r Response) GetStatus() int {
	reply := r.Ask(`kong.service.response.get_status`)
	status, _ := strconv.Atoi(reply)
	return status
}

func (r Response) GetHeaders(max_headers int) map[string]interface{} {
	var method string
	if max_headers == -1 {
		method = `kong.service.response.get_headers`
	} else {
		method = fmt.Sprintf(`kong.service.response.get_headers:%d`, max_headers)
	}

	headers := make(map[string]interface{})
	bridge.Unmarshal(r.Ask(method), &headers)
	return headers
}

func (r Response) GetHeader(name string) string {
	return r.Ask(fmt.Sprintf(`kong.service.response.get_header:%s`, name))
}
