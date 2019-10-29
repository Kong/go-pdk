package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var node Node
var ch chan string

func init() {
	ch = make(chan string)
	node = New(ch)
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
	assert.Equal(t, getName(func() { node.GetId() }), "kong.node.get_id:null")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := node.GetId(); res <- r }, "foo"), "foo")
	assert.Equal(t, getStrValue(func(res chan string) { r, _ := node.GetId(); res <- r }, ""), "")
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

	assert.Equal(t, "kong.node.get_memory_stats:null", getName(func() { node.GetMemoryStats() }))
	res := make(chan *MemoryStats)
	go func(res chan *MemoryStats) { r, _ := node.GetMemoryStats(); res <- r }(res)
	_ = <-ch
	ch <- stats
	statsO := <-res

	assert.Equal(t, 12288, statsO.LuaSharedDicts.Kong.AllocatedSlabs)
	assert.Equal(t, 24576, statsO.LuaSharedDicts.Kong.Capacity)
	assert.Equal(t, 12288, statsO.LuaSharedDicts.KongDbCache.AllocatedSlabs)
	assert.Equal(t, 12288, statsO.LuaSharedDicts.KongDbCache.Capacity)
	assert.Equal(t, 1102, statsO.WorkersLuaVms[0].HttpAllocatedGc)
	assert.Equal(t, 18004, statsO.WorkersLuaVms[0].Pid)
	assert.Equal(t, 1102, statsO.WorkersLuaVms[1].HttpAllocatedGc)
	assert.Equal(t, 18005, statsO.WorkersLuaVms[1].Pid)
}
