package log

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func mockLog(t *testing.T, s []bridgetest.MockStep) Log {
	return Log{bridge.New(bridgetest.Mock(t, s))}
}

func TestMessages(t *testing.T) {
	v, err := structpb.NewList([]interface{}{"Alo"})
	assert.NoError(t, err)

	log := mockLog(t, []bridgetest.MockStep{
		{Method: "kong.log.alert", Args: v, Ret: nil},
		{Method: "kong.log.crit", Args: v, Ret: nil},
		{Method: "kong.log.err", Args: v, Ret: nil},
		{Method: "kong.log.warn", Args: v, Ret: nil},
		{Method: "kong.log.notice", Args: v, Ret: nil},
		{Method: "kong.log.info", Args: v, Ret: nil},
		{Method: "kong.log.debug", Args: v, Ret: nil},
	})

	assert.NoError(t, log.Alert("Alo"))
	assert.NoError(t, log.Crit("Alo"))
	assert.NoError(t, log.Err("Alo"))
	assert.NoError(t, log.Warn("Alo"))
	assert.NoError(t, log.Notice("Alo"))
	assert.NoError(t, log.Info("Alo"))
	assert.NoError(t, log.Debug("Alo"))
}

func TestSerialize(t *testing.T) {
	log := mockLog(t, []bridgetest.MockStep{
		{Method: "kong.log.serialize", Args: nil, Ret: bridge.WrapString("{data...}")},
	})

	ret, err := log.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, "{data...}", ret)
}
