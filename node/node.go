package node

import (
	"encoding/json"
)

type Node struct {
	ch chan string
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

func NewNode(ch chan string) *Node {
	return &Node{ch: ch}
}

func (n *Node) GetId() string {
	n.ch <- `kong.node.get_id`
	return <-n.ch
}

func (n *Node) GetMemoryStats() *MemoryStats {
	n.ch <- `kong.node.get_memory_stats`
	stats := <-n.ch

	statsO := MemoryStats{}
	json.Unmarshal([]byte(stats), &statsO)
	return &statsO
}
