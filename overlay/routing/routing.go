package routing

import (
	"github.com/nknorg/nnet/common"
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/node"
)

// Router is an abstract routing layer that determines how to route a msg
type Router interface {
	GetNodeToRoute(*node.RemoteMessage) (*node.LocalNode, []*node.RemoteNode, error)
	Start() error
}

// Routing is the base struct for all routing
type Routing struct {
	localMsgChan chan *node.RemoteMessage
	rxMsgChan    chan *node.RemoteMessage
	common.LifeCycle
}

// NewRouting creates a new routing
func NewRouting(localMsgChan, rxMsgChan chan *node.RemoteMessage) (*Routing, error) {
	r := &Routing{
		localMsgChan: localMsgChan,
		rxMsgChan:    rxMsgChan,
	}
	return r, nil
}

// Start starts the msg handling process
func (r *Routing) Start(router Router, numWorkers int) error {
	r.StartOnce.Do(func() {
		for i := 0; i < numWorkers; i++ {
			go r.handleMsg(router)
		}
	})

	return nil
}

func (r *Routing) handleMsg(router Router) {
	var remoteMsg *node.RemoteMessage
	var localNode *node.LocalNode
	var remoteNode *node.RemoteNode
	var remoteNodes []*node.RemoteNode
	var err error

	for {
		if r.IsStopped() {
			return
		}

		select {
		case remoteMsg = <-r.rxMsgChan:
			localNode, remoteNodes, err = router.GetNodeToRoute(remoteMsg)
			if err != nil {
				log.Warn("Get next hop error:", err)
				continue
			}

			if localNode != nil {
				if len(remoteMsg.Msg.ReplyToId) > 0 {
					replyChan, ok := localNode.GetReplyChan(remoteMsg.Msg.ReplyToId)
					if ok && replyChan != nil {
						select {
						case replyChan <- remoteMsg:
						default:
							log.Warn("Reply chan full, discarding msg")
						}
					}
					continue
				}

				select {
				case r.localMsgChan <- remoteMsg:
				default:
					log.Warn("Router local msg chan full, discarding msg")
				}
			}

			for _, remoteNode = range remoteNodes {
				err = remoteNode.SendMessageAsync(remoteMsg.Msg)
				if err != nil {
					log.Warn("Send message to remote node error:", err)
				}
			}
		}
	}
}
