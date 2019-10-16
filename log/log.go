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

func (r *Log) Alert(args ...string) {
	r.ch <- fmt.Sprintf(`kong.log.alert:%s`, strings.Join(args, ""))
	_ = <-r.ch
}

func (r *Log) Crit(args ...string) {
	r.ch <- fmt.Sprintf(`kong.log.crit:%s`, strings.Join(args, ""))
	_ = <-r.ch
}

func (r *Log) Err(args ...string) {
	r.ch <- fmt.Sprintf(`kong.log.err:%s`, strings.Join(args, ""))
	_ = <-r.ch
}

func (r *Log) Warn(args ...string) {
	r.ch <- fmt.Sprintf(`kong.log.warn:%s`, strings.Join(args, ""))
	_ = <-r.ch
}

func (r *Log) Notice(args ...string) {
	r.ch <- fmt.Sprintf(`kong.log.notice:%s`, strings.Join(args, ""))
	_ = <-r.ch
}

func (r *Log) Info(args ...string) {
	r.ch <- fmt.Sprintf(`kong.log.info:%s`, strings.Join(args, ""))
	_ = <-r.ch
}
func (r *Log) Debug(args ...string) {
	r.ch <- fmt.Sprintf(`kong.log.debug:%s`, strings.Join(args, ""))
	_ = <-r.ch
}

func (r *Log) Serialize() string {
	r.ch <- fmt.Sprintf(`kong.log.serialize`)
	return <-r.ch
}
