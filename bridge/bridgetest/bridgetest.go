package bridgetest

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"testing"

	"google.golang.org/protobuf/proto"
)

type MockStep struct {
	Method string
	Args   proto.Message
	Ret    proto.Message
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

func Mock(t *testing.T, s []MockStep) net.Conn {
	conA, conB := net.Pipe()

	go func() {
		for i, stp := range s {
			d, err := readPbFrame(conB)
			if err != nil {
				t.Errorf("step %d readPbFrame(method): %s", i, err)
				break
			}
			if !bytes.Equal([]byte(stp.Method), d) {
				t.Errorf("step %d, expected method %v, found %v", i, []byte(stp.Method), d)
				break
			}

			d, err = readPbFrame(conB)
			if err != nil {
				t.Errorf("step %d, readPbFrame(args): %s", i, err)
				break
			}

			if stp.Args != nil {
				args_d, err := proto.Marshal(stp.Args)
				if err != nil {
					t.Errorf("step %d, Marshal(args): %s", i, err)
					break
				}

				if !bytes.Equal(args_d, d) {
					t.Errorf("step %d, expected %v(%v), received %v", i, stp.Args, args_d, d)
					break
				}
			}

			if stp.Ret != nil {
				ret_enc, err := proto.Marshal(stp.Ret)
				if err != nil {
					t.Errorf("step %d, Marshal(ret): %s", i, err)
					break
				}

				err = writePbFrame(conB, ret_enc)
				if err != nil {
					t.Errorf("step %d, writePbFrame(ret): %s", i, err)
					break
				}
			} else {
				err = writePbFrame(conB, []byte{})
				if err != nil {
					t.Errorf("step %d, writePbFrame(ret): %s", i, err)
					break
				}
			}
		}
		conB.Close()
	}()
	return conA
}

type mockEnvironment interface {
	Handle(method string, args_d []byte) []byte
	Errorf(format string, args ...interface{})
	IsRunning() bool
	SubscribeStatusChange(ch chan<- string)
}

func MockFunc(e mockEnvironment) net.Conn {
	conA, conB := net.Pipe()

	statusCh := make(chan string, 1)
	e.SubscribeStatusChange(statusCh)

	go func() {
		for {
			d, err := readPbFrame(conB)
			if err != nil {
				e.Errorf("Can't read method name")
				break
			}
			method := string(d)

			d, err = readPbFrame(conB)
			if err != nil {
				e.Errorf("Can't read method \"%v\" arguments", method)
				break
			}

			d = e.Handle(method, d)

			err = writePbFrame(conB, d)
			if err != nil {
				e.Errorf("Can't write back return values")
				break
			}

			select {
			case msg := <-statusCh:
				if msg == "finished" {
					return
				}
			default: // do nothing
			}
		}
	}()
	return conA
}
