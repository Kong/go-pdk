package server

import (
	"github.com/Kong/go-pdk/client"
	"github.com/Kong/go-pdk/entities"
	"github.com/Kong/go-pdk/node"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type StepErrorData struct {
	EventId int
	Data    Error
}

func (rh *rpcHandler) StepError(in StepErrorData, out *StepData) error {
	return rh.Step(StepData{
		EventId: in.EventId,
		Data:    in.Data,
	}, out)
}

type StepMultiMapData struct {
	EventId int
	Data    map[string][]string
}

func (rh *rpcHandler) StepMultiMap(in StepMultiMapData, out *StepData) error {
	return rh.Step(StepData{
		EventId: in.EventId,
		Data:    in.Data,
	}, out)
}

type StepCredentialData struct {
	EventId int
	Data    client.AuthenticatedCredential
}

func (rh *rpcHandler) StepCredential(in StepCredentialData, out *StepData) error {
	return rh.Step(StepData{
		EventId: in.EventId,
		Data:    in.Data,
	}, out)
}

type StepRouteData struct {
	EventId int
	Data    entities.Route
}

func (rh *rpcHandler) StepRoute(in StepRouteData, out *StepData) error {
	return rh.Step(StepData{
		EventId: in.EventId,
		Data:    in.Data,
	}, out)
}

type StepServiceData struct {
	EventId int
	Data    entities.Service
}

func (rh *rpcHandler) StepService(in StepServiceData, out *StepData) error {
	return rh.Step(StepData{
		EventId: in.EventId,
		Data:    in.Data,
	}, out)
}

type StepConsumerData struct {
	EventId int
	Data    entities.Consumer
}

func (rh *rpcHandler) StepConsumer(in StepConsumerData, out *StepData) error {
	return rh.Step(StepData{
		EventId: in.EventId,
		Data:    in.Data,
	}, out)
}

type StepMemoryStatsData struct {
	EventId int
	Data    node.MemoryStats
}

func (rh *rpcHandler) StepMemoryStats(in StepMemoryStatsData, out *StepData) error {
	return rh.Step(StepData{
		EventId: in.EventId,
		Data:    in.Data,
	}, out)
}
