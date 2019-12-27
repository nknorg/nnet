package chord

import (
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/overlay/routing"
)

const (
	// RelayRoutingNumWorkers determines how many concurrent goroutines are
	// handling relay messages
	RelayRoutingNumWorkers = 1
)

// RelayRouting is for message from a remote node to another remote node or
// local node
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

	rr := &RelayRouting{
		Routing: r,
		chord:   chord,
	}

	return rr, nil
}

// Start starts handling relay message from rxChan
func (rr *RelayRouting) Start() error {
	return rr.Routing.Start(rr, RelayRoutingNumWorkers)
}

// GetNodeToRoute returns the local node and remote nodes to route message to
func (rr *RelayRouting) GetNodeToRoute(remoteMsg *node.RemoteMessage) (*node.LocalNode, []*node.RemoteNode, error) {
	succ := rr.chord.successors.GetFirst()
	if succ == nil || BetweenLeftIncl(rr.chord.LocalNode.Id, succ.Id, remoteMsg.Msg.DestId) {
		return rr.chord.LocalNode, nil, nil
	}

	successors := rr.chord.successors.ToRemoteNodeList(true)
	for i := 0; i < len(successors)-1; i++ {
		if successors[i] == remoteMsg.RemoteNode {
			continue
		}
		if BetweenLeftIncl(successors[i].Id, successors[i+1].Id, remoteMsg.Msg.DestId) {
			return nil, []*node.RemoteNode{successors[i]}, nil
		}
	}

	for i := len(rr.chord.fingerTable) - 1; i >= 0; i-- {
		finger := rr.chord.fingerTable[i]
		first := finger.GetFirst()
		if first == nil {
			continue
		}
		if !BetweenIncl(rr.chord.LocalNode.Id, remoteMsg.Msg.DestId, first.Id) {
			continue
		}

		nextHop := first
		minRoundTripTime := first.GetRoundTripTime()
		for _, rn := range finger.ToRemoteNodeList(true) {
			if rn.IsOutbound && BetweenIncl(rr.chord.LocalNode.Id, remoteMsg.Msg.DestId, rn.Id) {
				rtt := rn.GetRoundTripTime()
				if minRoundTripTime == 0 || (rtt > 0 && rtt <= minRoundTripTime) {
					nextHop = rn
					minRoundTripTime = rtt
				}
			}
		}

		return nil, []*node.RemoteNode{nextHop}, nil
	}

	return nil, []*node.RemoteNode{succ}, nil
}
