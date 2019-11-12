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

func (b PdkBridge) SendCall(method string, args []interface{}) error {
	argsJson, err := json.Marshal(args)
	if err != nil {
		return err
	}

	call := method + ":" + string(argsJson)
	b.ch <- call

	return nil
}

func (b PdkBridge) ReturnReply() (string, error) {
	reply := <-b.ch
	if reply == "null" {
		return "", errors.New("null response")
	}

	return reply, nil
}

func (b PdkBridge) Ask(method string, args ...interface{}) (string, error) {
	if err := b.SendCall(method, args); err != nil {
		return "", err
	}

	return b.ReturnReply()
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
