package node

type Node struct {
	ch chan string
}

func NewNode(ch chan string) *Node {
	return &Node{ch: ch}
}

func (n *Node) GetId() string {
	n.ch <- `kong.node.get_id`
	return <- n.ch
}

// TODO get_memory_stats
