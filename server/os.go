package server

import (
	"flag"
	"log"
	"net"
	"os"
	"path"

	"github.com/ugorji/go/codec"
)

var (
	kongPrefix = flag.String("kong-prefix", "/usr/local/kong", "Kong prefix path (specified by the -p argument commonly used in the kong cli)")
	dump       = flag.Bool("dump", false, "Dump info about plugins")
	help       = flag.Bool("help", false, "Show usage info")
)

func init() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(2)
	}
}

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

	pth = path.Join(*kongPrefix, name+".socket")
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

type serverInfo struct {
	Protocol   string
	SocketPath string
	Plugins    []pluginInfo
}

func dumpInfo(rh rpcHandler) {
	info, err := rh.getInfo()
	if err != nil {
		log.Printf("getting plugin info: %s", err)
		return
	}

	socketPath, err := getSocketPath()
	if err != nil {
		log.Printf("getting Socket path: %s", err)
		return
	}

	var handle codec.JsonHandle
	enc := codec.NewEncoder(os.Stdout, &handle)
	err = enc.Encode(serverInfo{
		Protocol:   "ProtoBuf:1",
		SocketPath: socketPath,
		Plugins:    []pluginInfo{info},
	})
	if err != nil {
		log.Printf("encoding plugin info: %s", err)
	}
	os.Stdout.WriteString("\n")
}

