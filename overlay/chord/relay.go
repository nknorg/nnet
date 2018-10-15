package chord

import (
	"errors"

	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/overlay/routing"
)

const (
	// RelayRoutingNumWorkers determines how many concurrent goroutines are
	// handling Relay messages
	RelayRoutingNumWorkers = 1
)

// RelayRouting is for msg from a remote node to another remote node or local node
type RelayRouting struct {
	*routing.Routing
	chord *Chord
}

// NewRelayRouting creates a new RelayRouting
func NewRelayRouting(localMsgChan chan<- *node.RemoteMessage, rxMsgChan <-chan *node.RemoteMessage, chord *Chord) (*RelayRouting, error) {
	r, err := routing.NewRouting(localMsgChan, rxMsgChan)
	if err != nil {
		return nil, err
	}

	dr := &RelayRouting{
		Routing: r,
		chord:   chord,
	}

	return dr, nil
}

// Start starts handling relay msg from rxChan
func (dr *RelayRouting) Start() error {
	return dr.Routing.Start(dr, RelayRoutingNumWorkers)
}

// GetNodeToRoute returns the local node and remote nodes to route msg to
func (dr *RelayRouting) GetNodeToRoute(remoteMsg *node.RemoteMessage) (*node.LocalNode, []*node.RemoteNode, error) {
	succ := dr.chord.successors.GetFirst()
	if succ == nil {
		return nil, nil, errors.New("Local node has no successor yet")
	}

	if betweenLeftIncl(dr.chord.LocalNode.Id, succ.Id, remoteMsg.Msg.DestId) {
		return dr.chord.LocalNode, nil, nil
	}

	nextHop := succ
	minDist := distance(succ.Id, remoteMsg.Msg.DestId, dr.chord.nodeIDBits)

	for _, remoteNode := range dr.chord.successors.ToRemoteNodeList(false) {
		if remoteNode == remoteMsg.RemoteNode {
			continue
		}
		dist := distance(remoteNode.Id, remoteMsg.Msg.DestId, dr.chord.nodeIDBits)
		if dist.Cmp(minDist) < 0 {
			nextHop = remoteNode
			minDist = dist
		}
	}

	for _, finger := range dr.chord.fingerTable {
		for _, remoteNode := range finger.ToRemoteNodeList(false) {
			if remoteNode == remoteMsg.RemoteNode {
				continue
			}
			dist := distance(remoteNode.Id, remoteMsg.Msg.DestId, dr.chord.nodeIDBits)
			if dist.Cmp(minDist) < 0 {
				nextHop = remoteNode
				minDist = dist
			}
		}
	}

	return nil, []*node.RemoteNode{nextHop}, nil
}
