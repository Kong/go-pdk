package node

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var node *Node
var ch chan string

func init() {
	ch = make(chan string)
	node = &Node{ch: ch}
}

func getName(f func()) string {
	go f()
	name := <-ch
	ch <- ""
	return name
}

func getStrValue(f func(res chan string), val string) string {
	res := make(chan string)
	go f(res)
	_ = <-ch
	ch <- val
	return <-res
}

func TestGetId(t *testing.T) {
	assert.Equal(t, getName(func() { node.GetId() }), "kong.node.get_id")
	assert.Equal(t, getStrValue(func(res chan string) { res <- node.GetId() }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { res <- node.GetId() }, ""), "")
}

func TestGetMemoryStats(t *testing.T) {
	stats := `
	{
		"lua_shared_dicts": {
			"kong": {
				"allocated_slabs": 12288,
				"capacity": 24576
			},
			"kong_db_cache": {
				"allocated_slabs": 12288,
				"capacity": 12288
			}
		},
		"workers_lua_vms": [
			{
				"http_allocated_gc": 1102,
				"pid": 18004
			},
			{
				"http_allocated_gc": 1102,
				"pid": 18005
			}
		]
	}`

	assert.Equal(t, getName(func() { node.GetMemoryStats() }), "kong.node.get_memory_stats")
	res := make(chan *MemoryStats)
	go func(res chan *MemoryStats) { res <- node.GetMemoryStats() }(res)
	_ = <-ch
	ch <- stats
	statsO := <-res

	assert.Equal(t, statsO.LuaSharedDicts.Kong.AllocatedSlabs, 12288)
	assert.Equal(t, statsO.LuaSharedDicts.Kong.Capacity, 24576)
	assert.Equal(t, statsO.LuaSharedDicts.KongDbCache.AllocatedSlabs, 12288)
	assert.Equal(t, statsO.LuaSharedDicts.KongDbCache.Capacity, 12288)
	assert.Equal(t, statsO.WorkersLuaVms[0].HttpAllocatedGc, 1102)
	assert.Equal(t, statsO.WorkersLuaVms[0].Pid, 18004)
	assert.Equal(t, statsO.WorkersLuaVms[1].HttpAllocatedGc, 1102)
	assert.Equal(t, statsO.WorkersLuaVms[1].Pid, 18005)
}
