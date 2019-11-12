package bridge

import (
	"encoding/json"
	"errors"
)

type PdkBridge struct {
	ch chan string
}

type Response struct {
	Res string
	Err string
}

func New(ch chan string) PdkBridge {
	return PdkBridge{ch: ch}
}

func (b PdkBridge) Ask(method string, args ...interface{}) (string, error) {
	if argsJson, err := json.Marshal(args); err != nil {
		return "", err
	} else {
		call := method + ":" + string(argsJson)
		b.ch <- call
		reply := <-b.ch
		res := Response{}

		if err := json.Unmarshal([]byte(reply), &res); err != nil {
			return "", err
		}
		if res.Res != "" {
			return res.Res, nil
		} else if res.Err != "" {
			return "", errors.New(res.Err)
		} else {
			return "", nil
		}
	}
}

func Marshal(v interface{}) (string, error) {
	if b, err := json.Marshal(v); err != nil {
		return "", err
	} else {
		return string(b), err
	}
}

func Unmarshal(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}
