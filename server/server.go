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
	kongPrefix     = flag.String("kong-prefix", "/usr/local/kong", "Kong prefix path (specified by the -p argument commonly used in the kong cli)")
	dump = flag.Bool("dump", false, "Dump info about plugins")
	help = flag.Bool("help", false, "Show usage info")
)

func getName() (name string, err error) {
	execPath, err := os.Executable()
	if err != nil {
		return
	}

	name = path.Base(execPath)
	return
}

func getSocketPath() (pth string, err error) {
	name, err := getName()
	if err != nil {
		return
	}

	pth = path.Join(*kongPrefix, name + ".socket")
	return
}

func openSocket() (listener net.Listener, err error) {
	socketPath, err := getSocketPath()
	if err != nil {
		return
	}

	err = os.Remove(socketPath)
	if err != nil && !os.IsNotExist(err) {
		log.Printf(`removing "%s": %s`, socketPath, err)
		return
	}

	listener, err = net.Listen("unix", socketPath)
	if err != nil {
		log.Printf(`listen("%s"): %s`, socketPath, err)
		return
	}

	log.Printf("Listening on socket: %s", socketPath)
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
	err = enc.Encode([1]pluginInfo{ info })
	if err != nil {
		log.Printf("encoding plugin info: %s", err)
	}
	os.Stdout.WriteString("\n")
}

// Start the embedded plugin server
// Handles CLI flags, and returns immediately if appropriate.
// Otherwise, returns only if the server is stopped.
func StartServer(constructor func() interface{}, version string, priority int) error {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(2)
	}

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
