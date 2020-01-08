/*
Used internally for the RPC protocol.
*/
package bridge

import (
	"errors"
)

type PdkBridge struct {
	ch chan interface{}
}

type StepData struct {
	Method string
	Args []interface{}
}

func New(ch chan interface{}) PdkBridge {
	return PdkBridge{ch: ch}
}

func (b PdkBridge) Ask(method string, args ...interface{}) (interface{}, error) {
	b.ch <- StepData{ method, args }

	reply := <-b.ch
	if reply == "null" {
		return nil, errors.New("null response")
	}

	err, ok := reply.(error)
	if ok {
		return nil, err
	}

	return reply, nil
}

func (b PdkBridge) AskInt(method string, args ...interface{}) (i int, err error) {
	val, err := b.Ask(method, args...)
	if err != nil {
		return
	}

	var ok bool
	if i, ok = val.(int); !ok {
		err = ReturnTypeError("integer")
	}
	return
}

func (b PdkBridge) AskString(method string, args ...interface{}) (s string, err error) {
	val, err := b.Ask(method, args...)
	if err != nil {
		return
	}

	var ok bool
	if s, ok = val.(string); !ok {
		err = ReturnTypeError("string")
	}
	return
}

func (b PdkBridge) AskMap(method string, args ...interface{}) (m map[string]interface{}, err error) {
	val, err := b.Ask(method, args...)
	if err != nil {
		return
	}

	var ok bool
	if m, ok = val.(map[string]interface{}); !ok {
		err = ReturnTypeError("map[string]interface{}")
	}
	return
}


func ReturnTypeError(expected string) error {
	return errors.New("expected type: " + expected)
}
