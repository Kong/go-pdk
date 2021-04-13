/*
Node-level utilities
*/
package node

import (
	"github.com/Kong/go-pdk/server/kong_plugin_protocol"
	"github.com/Kong/go-pdk/bridge"
)

// Holds this module's functions.  Accessible as `kong.Node`
type Node struct {
	bridge.PdkBridge
}

type MemoryStats struct {
	LuaSharedDicts struct {
		Kong struct {
			AllocatedSlabs int64 `json:"allocated_slabs"`
			Capacity       int64 `json:"capacity"`
		} `json:"kong"`
		KongDbCache struct {
			AllocatedSlabs int64 `json:"allocated_slabs"`
			Capacity       int64 `json:"capacity"`
		} `json:"kong_db_cache"`
	} `json:"lua_shared_dicts"`
	WorkersLuaVms []struct {
		HttpAllocatedGc int64 `json:"http_allocated_gc"`
		Pid             int64 `json:"pid"`
	} `json:"workers_lua_vms"`
}

// Called by the plugin server at initialization.
// func New(conn net.Conn) Node {
// 	return Node{bridge{conn}}
// }

// kong.Node.GetId() returns the v4 UUID used by this node to describe itself.
func (n Node) GetId() (string, error) {
	return n.AskString(`kong.node.get_id`, nil)
}

// kong.Node.GetMemoryStats() returns memory usage statistics about this node.
func (n Node) GetMemoryStats() (MemoryStats, error) {
	out := new(kong_plugin_protocol.MemoryStats)
	err := n.Ask(`kong.node.get_memory_stats`, nil, out)
	if err != nil {
		return MemoryStats{}, err
	}

	ms := MemoryStats{}
	ms.LuaSharedDicts.Kong.AllocatedSlabs = out.LuaSharedDicts.Kong.AllocatedSlabs
	ms.LuaSharedDicts.Kong.Capacity = out.LuaSharedDicts.Kong.Capacity
	ms.LuaSharedDicts.KongDbCache.AllocatedSlabs = out.LuaSharedDicts.KongDbCache.AllocatedSlabs
	ms.LuaSharedDicts.KongDbCache.Capacity = out.LuaSharedDicts.KongDbCache.Capacity

// 	for i, wlv := range out.WorkersLuaVms {
// 		// somewhow add elements to ms.WorkersLuaVms[]
// 	}

	return ms, nil
}

