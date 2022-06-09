package server

import (
	"fmt"
)

// Incoming data for a new event.
// TODO: add some relevant data to reduce number of callbacks.
type StartEventData struct {
	InstanceId int    // Instance ID to start the event
	EventName  string // event name (not handler method name)
	// ....
}

type eventData struct {
	ipc chan interface{} // communication channel (TODO: use decoded structs)
}

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
