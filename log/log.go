/*
Write to log file.
*/
package log

import (
	"github.com/Kong/go-pdk/bridge"
	"google.golang.org/protobuf/types/known/structpb"
)

// Holds this module's functions.  Accessible as `kong.Log`
type Log struct {
	bridge.PdkBridge
}

func (r Log) doLog(method string, args []interface{}) error {
	l, err := structpb.NewList(args)
	if err != nil {
		return err
	}

	return r.Ask(method, l, nil)
}

func (r Log) Alert(args ...interface{}) error {
	return r.doLog(`kong.log.alert`, args)
}

func (r Log) Crit(args ...interface{}) error {
	return r.doLog(`kong.log.crit`, args)
}

func (r Log) Err(args ...interface{}) error {
	return r.doLog(`kong.log.err`, args)
}

func (r Log) Warn(args ...interface{}) error {
	return r.doLog(`kong.log.warn`, args)
}

func (r Log) Notice(args ...interface{}) error {
	return r.doLog(`kong.log.notice`, args)
}

func (r Log) Info(args ...interface{}) error {
	return r.doLog(`kong.log.info`, args)
}

func (r Log) Debug(args ...interface{}) error {
	return r.doLog(`kong.log.debug`, args)
}

// var (
// 	modeSet = map[string]string{ "mode": "set" }
// 	modeAdd = map[string]string{ "mode": "add" }
// 	modeReplace = map[string]string{ "mode": "replace" }
// )
//
// func (r Log) SetSerializeValue(key string, v interface{}) error {
// 	_, err := r.Ask(`kong.log.set_serialize_value`, key, v, modeSet)
// 	return err
// }
//
// func (r Log) SetSerializeValueAdd(key string, v interface{}) error {
// 	_, err := r.Ask(`kong.log.set_serialize_value`, key, v, modeAdd)
// 	return err
// }
//
// func (r Log) SetSerializeValueReplace(key string, v interface{}) error {
// 	_, err := r.Ask(`kong.log.set_serialize_value`, key, v, modeReplace)
// 	return err
// }

func (r Log) Serialize() (s string, err error) {
	return r.AskString(`kong.log.serialize`, nil)
}
