package response

import (
	"strconv"

	"github.com/kong/go-pdk/bridge"
)

type Response struct {
	bridge.PdkBridge
}

func New(ch chan string) Response {
	return Response{bridge.New(ch)}
}

func (r Response) GetStatus() (int, error) {
	reply, err := r.Ask(`kong.service.response.get_status`)
	if err != nil {
		return 0, err
	}

	status, _ := strconv.Atoi(reply)
	return status, nil
}

func (r Response) GetHeaders(max_headers int) (map[string]interface{}, error) {
	var res string
	var err error
	if max_headers == -1 {
		res, err = r.Ask(`kong.service.response.get_headers`)
	} else {
		res, err = r.Ask(`kong.service.response.get_headers`, max_headers)
	}
	if err != nil {
		return nil, err
	}

	headers := make(map[string]interface{})
	bridge.Unmarshal(res, &headers)
	return headers, nil
}

func (r Response) GetHeader(name string) (string, error) {
	return r.Ask(`kong.service.response.get_header`, name)
}
