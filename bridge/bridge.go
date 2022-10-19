/*
Used internally for the RPC protocol.
*/
package bridge

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"

	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"google.golang.org/protobuf/proto"
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

	if len > 0 {
		_, err = conn.Write(data)
	}

	return
}

func WrapString(s string) *kong_plugin_protocol.String {
	return &kong_plugin_protocol.String{V: s}
}

func WrapByteString(s []byte) *kong_plugin_protocol.ByteString {
	return &kong_plugin_protocol.ByteString{V: s}
}

func WrapHeaders(h map[string][]string) (*structpb.Struct, error) {
	h2 := make(map[string]interface{}, len(h))
	for k, v := range h {
		l := make([]interface{}, len(v))
		for i, v2 := range v {
			l[i] = v2
		}
		h2[k] = l
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
		case []interface{}:
			m2[k] = make([]string, len(v2))
			for i, v3 := range v2 {
				if s, ok := v3.(string); ok {
					m2[k][i] = s
				}
			}
		default:
			log.Printf("unexpected type %T on header %s:%v", v2, k, v2)
		}
	}

	return m2
}

func (b PdkBridge) Ask(method string, args proto.Message, out proto.Message) error {
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

func (b PdkBridge) Close() error {
	return b.conn.Close()
}

func ReturnTypeError(expected string) error {
	return errors.New("expected type: " + expected)
}
