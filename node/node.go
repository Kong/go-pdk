/*
Node-level utilities
*/
package node

import (
	"github.com/Kong/go-pdk/bridge"
)

// Holds this module's functions.  Accessible as `kong.Node`
type Node struct {
	bridge.PdkBridge
}

type MemoryStats struct {
	LuaSharedDicts struct {
		Kong struct {
			AllocatedSlabs int `json:"allocated_slabs"`
			Capacity       int `json:"capacity"`
		} `json:"kong"`
		KongDbCache struct {
			AllocatedSlabs int `json:"allocated_slabs"`
			Capacity       int `json:"capacity"`
		} `json:"kong_db_cache"`
	} `json:"lua_shared_dicts"`
	WorkersLuaVms []struct {
		HttpAllocatedGc int `json:"http_allocated_gc"`
		Pid             int `json:"pid"`
	} `json:"workers_lua_vms"`
}

// Called by the plugin server at initialization.
func New(ch chan interface{}) Node {
	return Node{bridge.New(ch)}
}

// kong.Node.GetId() returns the v4 UUID used by this node to describe itself.
func (n Node) GetId() (string, error) {
	return n.AskString(`kong.node.get_id`)
}

// kong.Node.GetMemoryStats() returns memory usage statistics about this node.
func (n Node) GetMemoryStats() (ms MemoryStats, err error) {
	val, err := n.Ask(`kong.node.get_memory_stats`)
	if err != nil {
		return
	}

	var ok bool
	if ms, ok = val.(MemoryStats); !ok {
		err = bridge.ReturnTypeError("MemoryStats")
	}
	return
}

// kong.Node.SetCtxShared() sets a value in the `kong.ctx.shared` request context table.
func (n Node) SetCtxShared(k string, value interface{}) error {
	_, err := n.Ask(`kong.set_ctx_shared`, k, value)
	return err
}

// kong.Node.GetCtxSharedAny() returns a value from the `kong.ctx.shared` request context table.
func (n Node) GetCtxSharedAny(k string) (interface{}, error) {
	return n.Ask(`kong.get_ctx_shared`, k)
}

// kong.Node.GetCtxSharedString() returns a string value from the `kong.ctx.shared` request context table.
func (n Node) GetCtxSharedString(k string) (string, error) {
	return n.AskString(`kong.get_ctx_shared`, k)
}

// kong.Node.GetCtxSharedFloat() returns a float value from the `kong.ctx.shared` request context table.
func (n Node) GetCtxSharedFloat(k string) (float64, error) {
	return n.AskFloat(`kong.get_ctx_shared`, k)
}
