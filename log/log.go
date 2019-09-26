package log

import (
	"fmt"
	"strings"
)

type Log struct {
	ch chan string
}

func NewLog(ch chan string) *Log {
	return &Log{ch: ch}
}

func (r *Log) Err(args ...string) {
	r.ch <- fmt.Sprintf(`kong.log.err:%s`, strings.Join(args, ""))
	_ = <-r.ch
}

func (r *Log) Serialize() string {
	r.ch <- fmt.Sprintf(`kong.log.serialize`)
	return <-r.ch
}
