package bridge

import (
	"encoding/json"
	"errors"
)

type PdkBridge struct {
	ch chan string
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
		if reply := <-b.ch; reply == "null" {
			return "", errors.New("null response")
		} else {
			return reply, nil
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
