package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/golang/protobuf/proto"
)

func servePb(conn net.Conn, rh *rpcHandler) (err error) {
	for {
		d, err := readPbFrame(conn)
		if err != nil {
			break
		}

		rd, err := codecPb(rh, conn, d)
		if err != nil {
			break
		}

		err = writePbFrame(conn, rd)
		if err != nil {
			break
		}
	}

	conn.Close()
	if err != nil {
		log.Print(err)
	}

	return
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

func codecPb(rh *rpcHandler, conn net.Conn, data []byte) (retData []byte, err error) {
	var m kong_plugin_protocol.RpcCall
	err = proto.Unmarshal(data, &m)
	if err != nil {
		return
	}

	rm, err := handlePbCmd(rh, conn, m)
	if err != nil {
		return
	}

	if rm != nil {
		retData, err = proto.Marshal(rm)
	}

	return
}

func pbInstanceStatus(status InstanceStatus) *kong_plugin_protocol.RpcReturn_InstanceStatus {
	return &kong_plugin_protocol.RpcReturn_InstanceStatus{
		InstanceStatus: &kong_plugin_protocol.InstanceStatus{
			Name:       status.Name,
			InstanceId: int32(status.Id),
			StartedAt:  status.StartTime,
		},
	}
}

func handlePbCmd(rh *rpcHandler, conn net.Conn, m kong_plugin_protocol.RpcCall) (rm *kong_plugin_protocol.RpcReturn, err error) {
	switch c := m.Call.(type) {
	case *kong_plugin_protocol.RpcCall_CmdGetPluginNames:
		// 		log.Printf("GetPluginNames: %v", c)

	case *kong_plugin_protocol.RpcCall_CmdGetPluginInfo:
		// 		log.Printf("GetPluginInfo: %v", c)

	case *kong_plugin_protocol.RpcCall_CmdStartInstance:
		config := PluginConfig{
			Name:   c.CmdStartInstance.Name,
			Config: c.CmdStartInstance.Config,
		}
		var status InstanceStatus
		err = rh.StartInstance(config, &status)
		if err != nil {
			return
		}

		rm = &kong_plugin_protocol.RpcReturn{
			Sequence: m.Sequence,
			Return:   pbInstanceStatus(status),
		}

	case *kong_plugin_protocol.RpcCall_CmdGetInstanceStatus:
		var status InstanceStatus
		err = rh.InstanceStatus(int(c.CmdGetInstanceStatus.InstanceId), &status)
		if err != nil {
			return
		}

		rm = &kong_plugin_protocol.RpcReturn{
			Sequence: m.Sequence,
			Return:   pbInstanceStatus(status),
		}

	case *kong_plugin_protocol.RpcCall_CmdCloseInstance:
		var status InstanceStatus
		err = rh.CloseInstance(int(c.CmdCloseInstance.InstanceId), &status)
		if err != nil {
			return
		}

		rm = &kong_plugin_protocol.RpcReturn{
			Sequence: m.Sequence,
			Return:   pbInstanceStatus(status),
		}

	case *kong_plugin_protocol.RpcCall_CmdHandleEvent:
		err = handlePbEvent(rh, conn, c.CmdHandleEvent)
		rm = &kong_plugin_protocol.RpcReturn{
			Sequence: m.Sequence,
		}

	default:
		err = fmt.Errorf("RPC call has unexpected type %T", c)
	}

	return
}

func handlePbEvent(rh *rpcHandler, conn net.Conn, e *kong_plugin_protocol.CmdHandleEvent) error {
	rh.lock.RLock()
	instance, ok := rh.instances[int(e.InstanceId)]
	rh.lock.RUnlock()
	if !ok {
		return fmt.Errorf("no plugin instance %d", e.InstanceId)
	}

	h, ok := instance.handlers[e.EventName]
	if !ok {
		return fmt.Errorf("undefined method %s", e.EventName)
	}

	pdk := pdk.Init(conn)

	h(pdk)
	writePbFrame(conn, []byte{})

	return nil
}

// Start the embedded plugin server, ProtoBuf version.
// Handles CLI flags, and returns immediately if appropriate.
// Otherwise, returns only if the server is stopped.
func StartServer(constructor func() interface{}, version string, priority int) error {
	rh := newRpcHandler(constructor, version, priority)

	if *dump {
		dumpInfo(*rh)
		return nil
	}

	listener, err := openSocket()
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go servePb(conn, rh)
	}

	return nil
}
