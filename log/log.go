package log

import (
	"fmt"
	"strings"
	"github.com/kong/go-pdk/bridge"
)

type Log struct {
	bridge.PdkBridge
}

func New(ch chan string) Log {
	return Log{bridge.New(ch)}
}

func (r Log) Alert(args ...string) {
	_ = r.Ask(fmt.Sprintf(`kong.log.alert:%s`, strings.Join(args, "")))
}

func (r Log) Crit(args ...string) {
	_ = r.Ask(fmt.Sprintf(`kong.log.crit:%s`, strings.Join(args, "")))
}

func (r Log) Err(args ...string) {
	_ = r.Ask(fmt.Sprintf(`kong.log.err:%s`, strings.Join(args, "")))
}

func (r Log) Warn(args ...string) {
	_ = r.Ask(fmt.Sprintf(`kong.log.warn:%s`, strings.Join(args, "")))
}

func (r Log) Notice(args ...string) {
	_ = r.Ask(fmt.Sprintf(`kong.log.notice:%s`, strings.Join(args, "")))
}

func (r Log) Info(args ...string) {
	_ = r.Ask(fmt.Sprintf(`kong.log.info:%s`, strings.Join(args, "")))
}

func (r Log) Debug(args ...string) {
	_ = r.Ask(fmt.Sprintf(`kong.log.debug:%s`, strings.Join(args, "")))
}

func (r Log) Serialize() string {
	return r.Ask(fmt.Sprintf(`kong.log.serialize`))
}
