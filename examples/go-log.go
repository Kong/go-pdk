/*
	A reimplementation of Kong's file-log plugin in Go.

	Here we use goroutines to detach execution from the
	request/response cycle.

	Global variables are also long-lived, and use Mutex
	locks to protect maps from being modified by
	concurrent execution.
*/

package main

import (
	"log"
	"os"
	"sync"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

// Start the embedded server
// Note that the parameters are just values, there's no
// need to give them exported names as required by
// the go-pluginserver.  Even the `New()` constructor
// function can be written inline.  (if it's more readable
// or not is up for debate)
func main() {
	server.StartServer(func() interface{} {
		return &Config{}
	}, "0.3", 2)
}

type Config struct {
	Path   string
	Reopen bool
}

// Global variables have the same lifespan
// as the service process.  No need
// to use any persistence to pass them
// from one event to another.
// Still, note that the service can be
// killed and respawned at any moment.
// Use some kind of external storage
// to save any state that must survive.
var fileDescriptors map[string]*os.File
var channels map[string]chan []byte

// multiple events can be handled concurrently
// global mutable maps (like those above) should
// be protected with locks
var fdmap_lock sync.Mutex
var chmap_lock sync.Mutex

func (conf Config) Log(kong *pdk.PDK) {
	ch, is_new := getChannel(conf.Path)
	if is_new {
		// this is where it "escapes" the event cycle
		go conf.collect(ch)
	}

	msg, err := kong.Log.Serialize()
	if err != nil {
		log.Print(err.Error())
		return
	}

	ch <- []byte(msg)
}

// get the channel associated with
// the file in the config.
func getChannel(path string) (chan []byte, bool) {
	chmap_lock.Lock()
	defer chmap_lock.Unlock()

	var is_new = false

	if channels == nil {
		channels = make(map[string]chan []byte)
	}

	ch, ok := channels[path]
	if !ok {
		channels[path] = make(chan []byte)
		is_new = true
		ch = channels[path]
	}

	return ch, is_new
}

//
// code below this line works "outside" Kong events and should not use the PDK
// --------------------------------
//

// Log collection
//
// Note that the goroutine that calls this function will
// persist for several Kong events so it can't use
// any PDK function.  The easiest way to enforce this
// is to factor the 'long lived' code in a function
// like this and *not* pass the `kong` variable.
//
// Here the `conf` object is the one passed from
// the event that spawned the goroutine.  It will
// not be updated on new events.
func (conf Config) collect(ch chan []byte) {
	for {
		b := <-ch

		fd, ok := conf.getFileDesc()
		if !ok {
			return
		}

		fd.Write(b)
		fd.Write([]byte("\n"))
	}
}

// get an open descriptor for the logfile.
//
// this is called on each new event to give
// the opportunity to close and reopen if
// requested by the configuration.
func (conf Config) getFileDesc() (*os.File, bool) {
	fdmap_lock.Lock()
	defer fdmap_lock.Unlock()

	if fileDescriptors == nil {
		fileDescriptors = make(map[string]*os.File)
	}

	fd, ok := fileDescriptors[conf.Path]
	if ok {
		if conf.Reopen {
			fd.Close()
			delete(fileDescriptors, conf.Path)
			ok = false
		}
	}

	if !ok {
		var err error
		fd, err = os.OpenFile(conf.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Printf("failed to open the file: %s", err.Error())
			return nil, false
		}
		fileDescriptors[conf.Path] = fd
	}

	return fd, true
}
