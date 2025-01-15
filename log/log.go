/*
Write to log file.
*/
package log

import (
	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
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
// 	modeSet = map[string]string{"mode": "set"}
// 	modeAdd     = map[string]string{"mode": "add"}
// 	modeReplace = map[string]string{"mode": "replace"}
// )

func (r Log) SetSerializeValue(key string, v interface{}) error {
	value, err := structpb.NewValue(v)
	if err != nil {
		return err
	}
	arg := kong_plugin_protocol.KV{
		K: key,
		V: value,
	}
	err = r.Ask(`kong.log.set_serialize_value`, &arg, nil)
	return err
}

// func (r Log) SetSerializeValueAdd(key string, v interface{}) error {
// 	_, err := r.Ask(`kong.log.set_serialize_value`, key, v, modeAdd)
// 	return err
// }

// func (r Log) SetSerializeValueReplace(key string, v interface{}) error {
// 	_, err := r.Ask(`kong.log.set_serialize_value`, key, v, modeReplace)
// 	return err
// }

func (r Log) Serialize() (s string, err error) {
	return r.AskString(`kong.log.serialize`, nil)
}
