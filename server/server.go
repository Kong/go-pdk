package server

import (
	"flag"
	"github.com/ugorji/go/codec"
	"log"
	"net"
	"net/rpc"
	"os"
	"path"
	"reflect"
)

var (
	dump = flag.Bool("dump-all-plugin", false, "Dump info about plugins")
)

func getName() (name string, err error) {
	execPath, err := os.Executable()
	if err != nil {
		return
	}

	name = path.Base(execPath)
	return
}

func getSocketPath() (string, error) {
	name, err := getName()
	if err != nil {
		return "", err
	}

	return name + ".socket", nil
}

func openSocket() (listener net.Listener, err error) {
	path, err := getName()
	if err != nil {
		return
	}

	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		log.Printf(`removing "%s": %s`, path, err)
		return
	}

	listener, err = net.Listen("unix", path)
	if err != nil {
		log.Printf(`listen("%s"): %s`, path, err)
		return
	}

	return
}

func runServer(listener net.Listener) {
	var handle codec.MsgpackHandle
	handle.ReaderBufferSize = 4096
	handle.WriterBufferSize = 4096
	handle.RawToString = true
	handle.MapType = reflect.TypeOf(map[string]interface{}(nil))

	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		enc := codec.NewEncoder(conn, &handle)
		_ = enc.Encode([]interface{}{2, "serverPid", os.Getpid()})

		rpcCodec := codec.MsgpackSpecRpc.ServerCodec(conn, &handle)
		go rpc.ServeCodec(rpcCodec)
	}
}

func dumpInfo(rh rpcHandler) {
	info, err := rh.getInfo()
	if err != nil {
		log.Printf("getting plugin info: %s", err)
		return
	}

	var handle codec.JsonHandle
	enc := codec.NewEncoder(os.Stdout, &handle)
	err = enc.Encode(info)
	if err != nil {
		log.Printf("encoding plugin info: %s", err)
	}
}

func StartServer(constructor func() interface{}) error {
	rh := newRpcHandler(constructor)

	if *dump {
		dumpInfo(*rh)
		return nil
	}

	listener, err := openSocket()
	if err != nil {
		return err
	}

	rpc.RegisterName("plugin", rh)
	runServer(listener)

	return nil
}
