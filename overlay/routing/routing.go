package routing

import (
	"errors"

	"github.com/nknorg/nnet/common"
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/util"
)

// Router is an abstract routing layer that determines how to route a message
type Router interface {
	Start() error
	ApplyMiddleware(interface{}) error
	GetNodeToRoute(*node.RemoteMessage) (*node.LocalNode, []*node.RemoteNode, error)
	SendMessage(Router, *node.RemoteMessage, bool) (<-chan *node.RemoteMessage, bool, error)
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

// SendMessage sends msg to the best next hop, returns reply chan (nil if if
// hasReply is false), if send success (which is true if successfully send
// message to at least one next hop), and aggregated errors during message
// sending
func (r *Routing) SendMessage(router Router, remoteMsg *node.RemoteMessage, hasReply bool) (<-chan *node.RemoteMessage, bool, error) {
	var shouldCallNextMiddleware bool

	localNode, remoteNodes, err := router.GetNodeToRoute(remoteMsg)
	if err != nil {
		return nil, false, err
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
		return nil, false, nil
	}

	if localNode == nil && len(remoteNodes) == 0 {
		return nil, false, errors.New("No node to route")
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
			return nil, true, nil
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
			return nil, true, nil
		}

		select {
		case r.localMsgChan <- remoteMsg:
		default:
			log.Warn("Router local msg chan full, discarding msg")
		}
	}

	var replyChan <-chan *node.RemoteMessage
	errs := util.NewErrors()
	success := false

	for _, remoteNode := range remoteNodes {
		// If there are multiple next hop, we only grab the first reply channel
		// because all msg have the same ID and will be using the same reply channel
		if hasReply && replyChan == nil {
			replyChan, err = remoteNode.SendMessage(remoteMsg.Msg, true)
		} else {
			_, err = remoteNode.SendMessage(remoteMsg.Msg, false)
		}

		if err != nil {
			errs = append(errs, err)
		} else {
			success = true
		}
	}

	return replyChan, success, errs.Merged()
}

// handleMsg starts to read received message from rxMsgChan, compute the route
// and dispatch message to local or remote nodes
func (r *Routing) handleMsg(router Router) {
	var remoteMsg *node.RemoteMessage
	var shouldCallNextMiddleware bool
	var err error

	for {
		if r.IsStopped() {
			return
		}

		remoteMsg = <-r.rxMsgChan

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

		_, _, err = r.SendMessage(router, remoteMsg, false)
		if err != nil {
			log.Warn(err)
		}
	}
}
