package node

import (
	"testing"

	"github.com/Kong/go-pdk/bridge"
	"github.com/Kong/go-pdk/bridge/bridgetest"
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/stretchr/testify/assert"
)

func mockNode(t *testing.T, s []bridgetest.MockStep) Node {
	return Node{bridge.New(bridgetest.Mock(t, s))}
}

func TestGetId(t *testing.T) {
	node := mockNode(t, []bridgetest.MockStep{
		{"kong.node.get_id", nil, bridge.WrapString("001:002:0003")},
	})

	ret, err := node.GetId()
	assert.NoError(t, err)
	assert.Equal(t, "001:002:0003", ret)
}

func TestGetMemoryStats(t *testing.T) {
	node := mockNode(t, []bridgetest.MockStep{
		{"kong.node.get_memory_stats", nil,
			&kong_plugin_protocol.MemoryStats{
				LuaSharedDicts: &kong_plugin_protocol.MemoryStats_LuaSharedDicts{
					Kong: &kong_plugin_protocol.MemoryStats_LuaSharedDicts_DictStats{
						AllocatedSlabs: 1027,
						Capacity:       4423543,
					},
					KongDbCache: &kong_plugin_protocol.MemoryStats_LuaSharedDicts_DictStats{
						AllocatedSlabs: 4093,
						Capacity:       3424875,
					},
				},
				WorkersLuaVms: []*kong_plugin_protocol.MemoryStats_WorkerLuaVm{
					{HttpAllocatedGc: 123456, Pid: 543},
					{HttpAllocatedGc: 345678, Pid: 876},
				},
			},
		},
	})

	ret, err := node.GetMemoryStats()
	assert.NoError(t, err)
	assert.Equal(t, MemoryStats{
		LuaSharedDicts: struct {
			Kong struct {
				AllocatedSlabs int64 "json:\"allocated_slabs\""
				Capacity       int64 "json:\"capacity\""
			} "json:\"kong\""
			KongDbCache struct {
				AllocatedSlabs int64 "json:\"allocated_slabs\""
				Capacity       int64 "json:\"capacity\""
			} "json:\"kong_db_cache\""
		}{
			Kong: struct {
				AllocatedSlabs int64 "json:\"allocated_slabs\""
				Capacity       int64 "json:\"capacity\""
			}{AllocatedSlabs: 1027, Capacity: 4423543},
			KongDbCache: struct {
				AllocatedSlabs int64 "json:\"allocated_slabs\""
				Capacity       int64 "json:\"capacity\""
			}{AllocatedSlabs: 4093, Capacity: 3424875},
		},
		WorkersLuaVms: []workerLuaVmStats{
			workerLuaVmStats{HttpAllocatedGc: 123456, Pid: 543},
			workerLuaVmStats{HttpAllocatedGc: 345678, Pid: 876},
		},
	}, ret)
}
