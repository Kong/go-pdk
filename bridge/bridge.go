package bridge

import "encoding/json"

type PdkBridge struct {
	ch chan string
}

func New(ch chan string) PdkBridge {
	return PdkBridge{ch: ch}
}

func (b PdkBridge) Ask (s string) string {
	b.ch <- s
	return <- b.ch
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
