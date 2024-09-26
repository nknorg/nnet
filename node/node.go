package node

import (
	"fmt"
	pbnode "github.com/nknorg/nnet/protobuf/node"
	"sync"

	"github.com/nknorg/nnet/common"
)

// Node is a remote or local node
type Node struct {
	sync.RWMutex
	*pbnode.Node
	common.LifeCycle
}

func newNode(n *pbnode.Node) (*Node, error) {
	node := &Node{
		Node: n,
	}
	return node, nil
}

// NewNode creates a node
func NewNode(id []byte, addr string) (*Node, error) {
	n := &pbnode.Node{
		Id:   id,
		Addr: addr,
	}
	return newNode(n)
}

func (n *Node) String() string {
	return fmt.Sprintf("%x@%s", n.Id, n.Addr)
}
