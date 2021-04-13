package server

import (
	"fmt"
	"github.com/Kong/go-pdk"
)

// Incoming data for a new event.
// TODO: add some relevant data to reduce number of callbacks.
type StartEventData struct {
	InstanceId int    // Instance ID to start the event
	EventName  string // event name (not handler method name)
	// ....
}

type eventData struct {
	id       int              // event id
	instance *instanceData    // plugin instance
	ipc      chan interface{} // communication channel (TODO: use decoded structs)
	pdk      *pdk.PDK         // go-pdk instance
}

func (rh *rpcHandler) addEvent(event *eventData) {
	rh.lock.Lock()
	defer rh.lock.Unlock()

	event.id = rh.nextEventId
	rh.nextEventId++
	rh.events[event.id] = event
}

// HandleEvent starts the call/{callback/response}*/finish cycle.
// More than one event can be run concurrenty for a single plugin instance,
// they all receive the same object instance, so should be careful if it's
// mutated or holds references to mutable data.
//
// RPC exported method
// func (rh *rpcHandler) HandleEvent(in StartEventData, out *StepData) error {
// 	rh.lock.RLock()
// 	instance, ok := rh.instances[in.InstanceId]
// 	rh.lock.RUnlock()
// 	if !ok {
// 		return fmt.Errorf("no plugin instance %d", in.InstanceId)
// 	}
//
// 	h, ok := instance.handlers[in.EventName]
// 	if !ok {
// 		return fmt.Errorf("undefined method %s", in.EventName)
// 	}
//
// 	ipc := make(chan interface{})
//
// 	event := eventData{
// 		instance: instance,
// 		ipc:      ipc,
// 		pdk:      pdk.Init(ipc),
// 	}
//
// 	rh.addEvent(&event)
//
// 	//log.Printf("Will launch goroutine for key %d / operation %s\n", key, op)
// 	go func() {
// 		_ = <-ipc
// 		h(event.pdk)
//
// 		func() {
// 			defer func() { recover() }()
// 			ipc <- "ret"
// 		}()
//
// 		rh.lock.Lock()
// 		defer rh.lock.Unlock()
// 		event.instance.lastEventTime = time.Now()
// 		delete(rh.events, event.id)
// 	}()
//
// 	ipc <- "run" // kickstart the handler
//
// 	*out = StepData{EventId: event.id, Data: <-ipc}
// 	return nil
// }

// A callback's response/request.
type StepData struct {
	EventId int         // event cycle to which this belongs
	Data    interface{} // carried data
}

// Step carries a callback's answer back from Kong to the plugin,
// the return value is either a new callback request or a finish signal.
//
// RPC exported method
func (rh *rpcHandler) Step(in StepData, out *StepData) error {
	rh.lock.RLock()
	event, ok := rh.events[in.EventId]
	rh.lock.RUnlock()
	if !ok {
		return fmt.Errorf("no running event %d", in.EventId)
	}

	event.ipc <- in.Data
	outStr := <-event.ipc
	*out = StepData{EventId: in.EventId, Data: outStr}

	return nil
}
