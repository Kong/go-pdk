package thread

import (
	"errors"
	"net"
	"github.com/kong/go-pdk/bridge"
)

type Thread struct {
	bridge.PdkBridge

	signalSocket net.Listener
	signalConn net.Conn
}

func New(ch chan string) Thread {
	return Thread{bridge.New(ch), nil, nil}
}

func (t Thread) makeSignalSocket() (string, error) {
	if t.signalSocket != nil {
		return "", errors.New("already detached")
	}

	uaddr, err := net.ResolveUnixAddr("unix", "./signalSocket")
	if err != nil {
		return "", err
	}

	t.signalSocket, err = net.ListenUnix("unix", uaddr)
	if err != nil {
		return "", err
	}

	return "./signalSocket", nil
}

func (t Thread) closeSignalSocket() {
	if t.signalConn != nil {
		t.signalConn.Close()
		t.signalConn = nil
	}

	if t.signalSocket != nil {
		t.signalSocket.Close()
		t.signalSocket = nil
	}
}

func (t Thread) establishSignalSocket() error {
	if t.signalSocket == nil {
		return errors.New("not detaching")
	}

	var err error
	t.signalConn, err = t.signalSocket.Accept()
	if err != nil {
		return err
	}

	return nil
}

func (t Thread) sendSignal() error {
	if t.signalConn == nil {
		return errors.New("no signal connection")
	}

	t.signalConn.Write([]byte(`.`))
	return nil
}

func (t Thread) Spawn(f func()) (string, error) {
	signalSocket, err := t.makeSignalSocket()
	if err != nil {
		return "", err
	}
	defer t.closeSignalSocket()

	err = t.SendCall(`kong.thread.yield`, signalSocket)

	doneChn := make(chan bool)
	go func() {
		f()
		doneChn <- true
	}()

	t.establishSignalSocket()

	_ = <- doneChn
	t.sendSignal()
	return t.ReturnReply()
}
