package chord

import (
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/nknorg/nnet/config"
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/overlay"
	"github.com/nknorg/nnet/overlay/routing"
	"github.com/nknorg/nnet/protobuf"
)

const (
	// How many concurrent goroutines are handling messages
	numWorkers = 1
)

// Chord is the overlay network based on Chord DHT
type Chord struct {
	*overlay.Overlay
	nodeIDBits           uint32
	minStabilizeInterval time.Duration
	maxStabilizeInterval time.Duration
	successors           *NeighborList
	predecessors         *NeighborList
	fingerTable          []*NeighborList
	neighbors            *NeighborList
	*middlewareStore
}

// NewChord creates a Chord overlay network
func NewChord(localNode *node.LocalNode, conf *config.Config) (*Chord, error) {
	ovl, err := overlay.NewOverlay(localNode)
	if err != nil {
		return nil, err
	}

	nodeIDBits := conf.NodeIDBytes * 8

	next := nextID(localNode.Id, nodeIDBits)
	prev := prevID(localNode.Id, nodeIDBits)

	successors, err := NewNeighborList(next, prev, nodeIDBits, conf.MinNumSuccessors, false)
	if err != nil {
		return nil, err
	}

	predecessors, err := NewNeighborList(prev, next, nodeIDBits, conf.MinNumPredecessors, true)
	if err != nil {
		return nil, err
	}

	fingerTable := make([]*NeighborList, nodeIDBits)
	for i := uint32(0); i < nodeIDBits; i++ {
		startID := powerOffset(localNode.Id, i, nodeIDBits)
		endID := prevID(powerOffset(localNode.Id, i+1, nodeIDBits), nodeIDBits)
		fingerTable[i], err = NewNeighborList(startID, endID, nodeIDBits, conf.NumFingerSuccessors, false)
		if err != nil {
			return nil, err
		}
	}

	neighbors, err := NewNeighborList(next, prev, nodeIDBits, 0, false)
	if err != nil {
		return nil, err
	}

	middlewareStore := newMiddlewareStore()

	c := &Chord{
		Overlay:              ovl,
		nodeIDBits:           nodeIDBits,
		minStabilizeInterval: conf.MinStabilizeInterval,
		maxStabilizeInterval: conf.MaxStabilizeInterval,
		successors:           successors,
		predecessors:         predecessors,
		fingerTable:          fingerTable,
		neighbors:            neighbors,
		middlewareStore:      middlewareStore,
	}

	directRxMsgChan, err := localNode.GetRxMsgChan(protobuf.DIRECT)
	if err != nil {
		return nil, err
	}
	directRouting, err := routing.NewDirectRouting(ovl.LocalMsgChan, directRxMsgChan)
	if err != nil {
		return nil, err
	}
	err = ovl.AddRouter(protobuf.DIRECT, directRouting)
	if err != nil {
		return nil, err
	}

	relayRxMsgChan, err := localNode.GetRxMsgChan(protobuf.RELAY)
	if err != nil {
		return nil, err
	}
	relayRouting, err := NewRelayRouting(ovl.LocalMsgChan, relayRxMsgChan, c)
	if err != nil {
		return nil, err
	}
	err = ovl.AddRouter(protobuf.RELAY, relayRouting)
	if err != nil {
		return nil, err
	}

	err = localNode.ApplyMiddleware(node.RemoteNodeReady(func(rn *node.RemoteNode) bool {
		c.addNeighbor(rn)
		return true
	}))
	if err != nil {
		return nil, err
	}

	err = localNode.ApplyMiddleware(node.RemoteNodeDisconnected(func(rn *node.RemoteNode) bool {
		c.removeNeighbor(rn)
		return true
	}))
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Start starts the runtime loop of the chord network
func (c *Chord) Start() error {
	c.StartOnce.Do(func() {
		var joinOnce sync.Once

		err := c.ApplyMiddleware(SuccessorAdded(func(remoteNode *node.RemoteNode, index int) bool {
			joinOnce.Do(func() {
				// prev is used to prevent msg being routed to self
				prev := prevID(c.LocalNode.Id, c.nodeIDBits)
				succs, err := c.FindSuccessors(prev, c.successors.Cap())
				if err != nil {
					log.Error("Join failed:", err)
				}

				for _, succ := range succs {
					if CompareID(succ.Id, c.LocalNode.Id) != 0 {
						err = c.Connect(succ.Addr, succ.Id)
						if err != nil {
							log.Error(err)
						}
					}
				}
			})
			return true
		}))
		if err != nil {
			c.Stop(err)
		}

		for i := 0; i < numWorkers; i++ {
			go c.handleMsg()
		}

		err = c.StartRouters()
		if err != nil {
			c.Stop(err)
		}

		go c.stabilize()
	})

	return nil
}

// Stop stops the chord network
func (c *Chord) Stop(err error) {
	c.StopOnce.Do(func() {
		if err != nil {
			log.Warnf("Chord overlay stops because of error: %s", err)
		} else {
			log.Infof("Chord overlay stops")
		}

		c.LifeCycle.Stop()
	})
}

// Join joins an existing chord network starting from the seedNodeAddr
func (c *Chord) Join(seedNodeAddr string) error {
	err := c.Connect(seedNodeAddr, nil)
	if err != nil {
		return err
	}

	return nil
}

// handleMsg starts a loop that handles received msg
func (c *Chord) handleMsg() {
	var remoteMsg *node.RemoteMessage
	var shouldLocalNodeHandleMsg bool
	var err error

	for {
		if c.IsStopped() {
			return
		}

		select {
		case remoteMsg = <-c.LocalMsgChan:
			shouldLocalNodeHandleMsg, err = c.handleRemoteMessage(remoteMsg)
			if err != nil {
				log.Error(err)
				continue
			}

			if shouldLocalNodeHandleMsg {
				err = c.LocalNode.HandleRemoteMessage(remoteMsg)
				if err != nil {
					log.Error(err)
					continue
				}
			}
		}
	}
}

// stabilize periodically updates successors and fingerTable to keep topology
// correct
func (c *Chord) stabilize() {
	go c.updateSuccAndPred()
	go c.updateNonEmptyFinger()
	go c.updateEmptyFinger()
}

// updateSuccAndPred periodically updates successors and predecessors
func (c *Chord) updateSuccAndPred() {
	var err error

	for {
		if c.IsStopped() {
			return
		}

		time.Sleep(randDuration(c.minStabilizeInterval, c.maxStabilizeInterval))

		err = c.updateNeighborList(c.successors)
		if err != nil {
			log.Error("Update successors error:", err)
		}

		err = c.updateNeighborList(c.predecessors)
		if err != nil {
			log.Error("Update predecessor error:", err)
		}
	}
}

// updateSuccAndPred periodically updates non-empty finger table items
func (c *Chord) updateNonEmptyFinger() {
	var err error
	var finger *NeighborList

	for {
		for _, finger = range c.fingerTable {
			if finger.IsEmpty() {
				continue
			}

			if c.IsStopped() {
				return
			}

			time.Sleep(randDuration(c.minStabilizeInterval, c.maxStabilizeInterval))

			err = c.updateNeighborList(finger)
			if err != nil {
				log.Error("Update finger table error:", err)
			}
		}

		// to prevent endless looping when fingerTable is all empty
		time.Sleep(randDuration(c.minStabilizeInterval, c.maxStabilizeInterval))
	}
}

// updateSuccAndPred periodically updates empty finger table items
func (c *Chord) updateEmptyFinger() {
	var err error
	var i int
	var succs []*protobuf.Node

	for {
		for i = 0; i < len(c.fingerTable); i++ {
			if !c.fingerTable[i].IsEmpty() {
				continue
			}

			if c.IsStopped() {
				return
			}

			time.Sleep(randDuration(c.minStabilizeInterval, c.maxStabilizeInterval))

			succs, err = c.FindSuccessors(c.fingerTable[i].startID, 1)
			if err != nil {
				// log.Error("Find successor for finger table error:", err)
				continue
			}

			if succs == nil || len(succs) == 0 {
				continue
			}

			for i < len(c.fingerTable) {
				if betweenIncl(c.fingerTable[i].startID, c.fingerTable[i].endID, succs[0].Id) {
					existing := c.fingerTable[i].GetFirst()
					if existing == nil || betweenLeftIncl(c.fingerTable[i].startID, existing.Id, succs[0].Id) {
						// TODO: check if in connecting list
						err = c.Connect(succs[0].Addr, succs[0].Id)
						if err != nil {
							log.Error("Connect to new node error:", err)
						}
					}
					break
				}
				i++
			}
		}

		// to prevent endless looping when fingerTable is all non-empty
		time.Sleep(randDuration(c.minStabilizeInterval, c.maxStabilizeInterval))
	}
}

// GetSuccAndPred sends a GetSuccAndPred message to remote node and returns its
// successors and predecessor if no error occured
func GetSuccAndPred(rn *node.RemoteNode, numSucc, numPred uint32) ([]*protobuf.Node, []*protobuf.Node, error) {
	msg, err := NewGetSuccAndPredMessage(numSucc, numPred)
	if err != nil {
		return nil, nil, err
	}

	reply, err := rn.SendMessageSync(msg)
	if err != nil {
		return nil, nil, err
	}

	replyBody := &protobuf.GetSuccAndPredReply{}
	err = proto.Unmarshal(reply.Msg.Message, replyBody)
	if err != nil {
		return nil, nil, err
	}

	return replyBody.Successors, replyBody.Predecessors, nil
}

// FindSuccessors sends a FindSuccessors message and returns numSucc successors
// of a given key id
func (c *Chord) FindSuccessors(key []byte, numSucc uint32) ([]*protobuf.Node, error) {
	succ := c.successors.GetFirst()
	if succ != nil && between(c.LocalNode.Id, succ.Id, key) {
		succs := c.successors.ToProtoNodeList(true)
		if succs != nil && uint32(len(succs)) > numSucc {
			succs = succs[:numSucc]
		}
		return succs, nil
	}

	msg, err := NewFindSuccessorsMessage(key, numSucc)
	if err != nil {
		return nil, err
	}

	reply, success, err := c.SendMessageSync(msg, protobuf.RELAY)
	if !success {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	replyBody := &protobuf.FindSuccessorsReply{}
	err = proto.Unmarshal(reply.Message, replyBody)
	if err != nil {
		return nil, err
	}

	return replyBody.Successors, nil
}
