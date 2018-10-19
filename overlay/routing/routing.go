package routing

import (
	"github.com/nknorg/nnet/common"
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/node"
)

// Router is an abstract routing layer that determines how to route a message
type Router interface {
	Start() error
	ApplyMiddleware(interface{}) error
	GetNodeToRoute(*node.RemoteMessage) (*node.LocalNode, []*node.RemoteNode, error)
}

// Routing is the base struct for all routing
type Routing struct {
	localMsgChan chan<- *node.RemoteMessage
	rxMsgChan    <-chan *node.RemoteMessage
	*middlewareStore
	common.LifeCycle
}

// NewRouting creates a new routing
func NewRouting(localMsgChan chan<- *node.RemoteMessage, rxMsgChan <-chan *node.RemoteMessage) (*Routing, error) {
	r := &Routing{
		localMsgChan:    localMsgChan,
		rxMsgChan:       rxMsgChan,
		middlewareStore: newMiddlewareStore(),
	}
	return r, nil
}

// Start starts the message handling process
func (r *Routing) Start(router Router, numWorkers int) error {
	r.StartOnce.Do(func() {
		for i := 0; i < numWorkers; i++ {
			go r.handleMsg(router)
		}
	})

	return nil
}

// handleMsg starts to read received message from rxMsgChan, compute the route
// and dispatch message to local or remote nodes
func (r *Routing) handleMsg(router Router) {
	var remoteMsg *node.RemoteMessage
	var localNode *node.LocalNode
	var remoteNode *node.RemoteNode
	var remoteNodes []*node.RemoteNode
	var shouldCallNextMiddleware bool
	var err error

	for {
		if r.IsStopped() {
			return
		}

		select {
		case remoteMsg = <-r.rxMsgChan:
			r.middlewareStore.RLock()
			for _, f := range r.middlewareStore.remoteMessageArrived {
				remoteMsg, shouldCallNextMiddleware = f(remoteMsg)
				if remoteMsg == nil || !shouldCallNextMiddleware {
					break
				}
			}
			r.middlewareStore.RUnlock()

			if remoteMsg == nil {
				continue
			}

			localNode, remoteNodes, err = router.GetNodeToRoute(remoteMsg)
			if err != nil {
				log.Warn("Get next hop error:", err)
				continue
			}

			r.middlewareStore.RLock()
			for _, f := range r.middlewareStore.remoteMessageRouted {
				remoteMsg, localNode, remoteNodes, shouldCallNextMiddleware = f(remoteMsg, localNode, remoteNodes)
				if remoteMsg == nil || !shouldCallNextMiddleware {
					break
				}
			}
			r.middlewareStore.RUnlock()

			if remoteMsg == nil {
				continue
			}

			if localNode != nil {
				r.middlewareStore.RLock()
				for _, f := range r.middlewareStore.remoteMessageReceived {
					remoteMsg, shouldCallNextMiddleware = f(remoteMsg)
					if remoteMsg == nil || !shouldCallNextMiddleware {
						break
					}
				}
				r.middlewareStore.RUnlock()

				if remoteMsg == nil {
					continue
				}

				if len(remoteMsg.Msg.ReplyToId) > 0 {
					replyChan, ok := localNode.GetReplyChan(remoteMsg.Msg.ReplyToId)
					if ok && replyChan != nil {
						select {
						case replyChan <- remoteMsg:
						default:
							log.Warn("Reply chan unavailable or full, discarding msg")
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
