package node

import (
	"github.com/Kong/go-pdk/bridge"
)

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

func New(ch chan string) Node {
	return Node{bridge.New(ch)}
}

func (n Node) GetId() (string, error) {
	return n.Ask(`kong.node.get_id`)
}

func (n Node) GetMemoryStats() (*MemoryStats, error) {
	statsO := MemoryStats{}
	res, err := n.Ask(`kong.node.get_memory_stats`)
	if err != nil {
		return nil, err
	}
	bridge.Unmarshal(res, &statsO)
	return &statsO, nil
}
