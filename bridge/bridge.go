/*
Used internally for the RPC protocol.
*/
package bridge

import (
	"errors"
	"encoding/binary"
	"io"
	"log"
	"net"
	"github.com/golang/protobuf/proto"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"google.golang.org/protobuf/types/known/structpb"
)

type PdkBridge struct {
	conn net.Conn
}

type StepData struct {
	Method string
	Args   []interface{}
}

func New(conn net.Conn) PdkBridge {
	return PdkBridge{
		conn: conn,
	}
}

func readPbFrame(conn net.Conn) (data []byte, err error) {
	var len uint32
	err = binary.Read(conn, binary.LittleEndian, &len)
	if err != nil {
		return
	}

	data = make([]byte, len)
	if data == nil {
		return nil, errors.New("no memory")
	}

	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}

	return
}

func writePbFrame(conn net.Conn, data []byte) (err error) {
	var len uint32 = uint32(len(data))
	err = binary.Write(conn, binary.LittleEndian, len)
	if err != nil {
		return
	}

	_, err = conn.Write(data)

	return
}


func WrapString(s string) *kong_plugin_protocol.String {
	return &kong_plugin_protocol.String{V: s}
}

func WrapHeaders(h map[string][]string) (*structpb.Struct, error) {
	h2 := make(map[string]interface{}, len(h))
	for k, v := range h {
		h2[k] = v
	}

	st, err := structpb.NewStruct(h2)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func UnwrapHeaders(st *structpb.Struct) map[string][]string {
	m := st.AsMap()
	m2 := make(map[string][]string)
	for k, v := range m {
		switch v2 := v.(type) {
			case string:
				m2[k] = []string{v2}
			case []string:
				m2[k] = v2
			default:
				log.Printf("unexpected type %T on header %s:%v", v2, k, v2)
		}
	}

	return m2
}


func (b PdkBridge) Ask(method string, args proto.Message, out proto.Message) error {
// 	log.Printf("Ask: method: [%v], args: [%#v], out: [%T]", method, args, out)
	err := writePbFrame(b.conn, []byte(method))
	if err != nil {
		return err
	}

	var args_d []byte

	if args != nil {
		args_d, err = proto.Marshal(args)
		if err != nil {
			return err
		}
	}

	err = writePbFrame(b.conn, args_d)
	if err != nil {
		return err
	}

	out_d, err := readPbFrame(b.conn)
	if err != nil {
		return err
	}

	if out != nil {
		err = proto.Unmarshal(out_d, out)
	}

	return err
}


func (b PdkBridge) AskString(method string, args proto.Message) (string, error) {
	out := new(kong_plugin_protocol.String)
	err := b.Ask(method, args, out)
	return out.V, err
}

func (b PdkBridge) AskInt(method string, args proto.Message) (int, error) {
	out := new(kong_plugin_protocol.Int)
	err := b.Ask(method, args, out)
	return int(out.V), err
}

func (b PdkBridge) AskNumber(method string, args proto.Message) (float64, error) {
	out := new(kong_plugin_protocol.Number)
	err := b.Ask(method, args, out)
	return out.V, err
}

func (b PdkBridge) AskValue(method string, args proto.Message) (interface{}, error) {
	out := new(structpb.Value)
	err := b.Ask(method, args, out)
	if err != nil {
		return nil, err
	}

	return out.AsInterface(), nil
}

// func (b PdkBridge) Ask(method string, args ...interface{}) (interface{}, error) {
// 	b.ch <- StepData{method, args}
//
// 	reply := <-b.ch
//
// 	err, ok := reply.(error).Printf
// 	if ok {
// 		return nil, err
// 	}
//
// 	return reply, nil
// }
/*
func (b PdkBridge) AskClose(method string, args ...interface{}) {
	b.ch <- StepData{ method, args }
	close(b.ch)
}

func (b PdkBridge) AskInt(method string, args ...interface{}) (i int, err error) {
	val, err := b.Ask(method, args...)
	if err != nil {
		return
	}
	if val == nil {
		err = errors.New("null response")
		return
	}

	switch val := val.(type) {
	case int:
		i = int(val)
	case int8:
		i = int(val)
	case int16:
		i = int(val)
	case int32:
		i = int(val)
	case int64:
		i = int(val)
	case uint:
		i = int(val)
	case uint8:
		i = int(val)
	case uint16:
		i = int(val)
	case uint32:
		i = int(val)
	case uint64:
		i = int(val)
	default:
		err = ReturnTypeError("integer")
	}
	return
}

func (b PdkBridge) AskFloat(method string, args ...interface{}) (f float64, err error) {
	val, err := b.Ask(method, args...)
	if err != nil {
		return
	}
	if val == nil {
		err = errors.New("null response")
		return
	}

	switch val := val.(type) {
	case int:
		f = float64(val)
	case int8:
		f = float64(val)
	case int16:
		f = float64(val)
	case int32:
		f = float64(val)
	case int64:
		f = float64(val)
	case uint:
		f = float64(val)
	case uint8:
		f = float64(val)
	case uint16:
		f = float64(val)
	case uint32:
		f = float64(val)
	case uint64:
		f = float64(val)
	case float32:
		f = float64(val)
	case float64:
		f = float64(val)
	default:
		err = ReturnTypeError("float")
	}
	return
}

func (b PdkBridge) AskString(method string, args ...interface{}) (s string, err error) {
	val, err := b.Ask(method, args...)
	if err != nil {
		return
	}
	if val == nil {
		err = errors.New("null response")
		return
	}

	var ok bool
	if s, ok = val.(string); !ok {
		err = ReturnTypeError("string")
	}
	return
}

func (b PdkBridge) AskMap(method string, args ...interface{}) (m map[string][]string, err error) {
	val, err := b.Ask(method, args...)
	if err != nil {
		return
	}

	var ok bool
	if m, ok = val.(map[string][]string); !ok {
		err = ReturnTypeError("map[string][]string")
	}
	return
}*/

func (b PdkBridge) Close() error {
	return b.conn.Close()
}

func ReturnTypeError(expected string) error {
	return errors.New("expected type: " + expected)
}
