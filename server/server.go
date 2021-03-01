/*
Package Kong/go-pdk/server implements an embedded plugin server.

To use, add a main() function:

    func main () {
      server.StartServer(New, Version, Priority)
    }

and compile as an executable with the standard `go build` command.
*/
package server

import (
	"github.com/ugorji/go/codec"
	"net"
	"net/rpc"
	"os"
	"reflect"
)


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

// Start the embedded plugin server
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

	err = rpc.RegisterName("plugin", rh)
	if err != nil {
		return err
	}

	runServer(listener)

	return nil
}
