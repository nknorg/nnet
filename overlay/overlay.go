package overlay

import (
	"errors"
	"fmt"
	"time"

	"github.com/nknorg/nnet/common"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/overlay/routing"
	"github.com/nknorg/nnet/protobuf"
)

// Overlay is an abstract overlay network
type Overlay struct {
	LocalNode    *node.LocalNode
	LocalMsgChan chan *node.RemoteMessage
	routers      map[protobuf.RoutingType]routing.Router
	common.LifeCycle
}

// NewOverlay creates a new overlay network
func NewOverlay(localNode *node.LocalNode) (*Overlay, error) {
	if localNode == nil {
		return nil, errors.New("Local node is nil")
	}

	overlay := &Overlay{
		LocalNode:    localNode,
		LocalMsgChan: make(chan *node.RemoteMessage, localNode.OverlayLocalMsgChanLen),
		routers:      make(map[protobuf.RoutingType]routing.Router),
	}
	return overlay, nil
}

// GetLocalNode returns the local node
func (ovl *Overlay) GetLocalNode() *node.LocalNode {
	return ovl.LocalNode
}

// AddRouter adds a router for a routingType, and returns error if router has
// already benn added for the type
func (ovl *Overlay) AddRouter(routingType protobuf.RoutingType, router routing.Router) error {
	_, ok := ovl.routers[routingType]
	if ok {
		return fmt.Errorf("Router for type %v is already added", routingType)
	}
	ovl.routers[routingType] = router
	return nil
}

// GetRouter gets a router for a routingType, and returns error if router of the
// type has not benn added yet
func (ovl *Overlay) GetRouter(routingType protobuf.RoutingType) (routing.Router, error) {
	router, ok := ovl.routers[routingType]
	if !ok {
		return nil, fmt.Errorf("Router for type %v has not been added yet", routingType)
	}
	return router, nil
}

// GetRouters gets a list of routers added
func (ovl *Overlay) GetRouters() []routing.Router {
	routers := make([]routing.Router, 0)
	for _, router := range ovl.routers {
		routers = append(routers, router)
	}
	return routers
}

// SetRouter sets a router for a routingType regardless of whether a router has
// been set up for this type or not
func (ovl *Overlay) SetRouter(routingType protobuf.RoutingType, router routing.Router) {
	ovl.routers[routingType] = router
}

// StartRouters starts all routers added to overlay network
func (ovl *Overlay) StartRouters() error {
	for _, router := range ovl.routers {
		err := router.Start()
		if err != nil {
			return nil
		}
	}

	return nil
}

// StopRouters stops all routers added to overlay network
func (ovl *Overlay) StopRouters(err error) {
	for _, router := range ovl.routers {
		router.Stop(err)
	}
}

// SendMessage sends msg to the best next hop, returns reply chan (nil if if
// hasReply is false), if send success (which is true if successfully send
// message to at least one next hop), and aggregated errors during message
// sending.
func (ovl *Overlay) SendMessage(msg *protobuf.Message, routingType protobuf.RoutingType, hasReply bool, replyTimeout time.Duration) (<-chan *node.RemoteMessage, bool, error) {
	router, err := ovl.GetRouter(routingType)
	if err != nil {
		return nil, false, err
	}

	return router.SendMessage(router, &node.RemoteMessage{Msg: msg}, hasReply, replyTimeout)
}

// SendMessageAsync sends msg to the best next hop, returns if send success
// (which is true if successfully send message to at least one next hop), and
// aggretated error during message sending
func (ovl *Overlay) SendMessageAsync(msg *protobuf.Message, routingType protobuf.RoutingType) (bool, error) {
	_, success, err := ovl.SendMessage(msg, routingType, false, 0)
	return success, err
}

// SendMessageSync sends msg to the best next hop, returns reply message, if
// send success (which is true if successfully send message to at least one next
// hop), and aggregated error during message sending, will also returns error if
// haven't receive reply within replyTimeout. Will use default reply timeout if
// replyTimeout = 0.
func (ovl *Overlay) SendMessageSync(msg *protobuf.Message, routingType protobuf.RoutingType, replyTimeout time.Duration) (*protobuf.Message, bool, error) {
	if replyTimeout == 0 {
		replyTimeout = ovl.LocalNode.DefaultReplyTimeout
	}

	replyChan, success, err := ovl.SendMessage(msg, routingType, true, replyTimeout)
	if !success {
		return nil, false, err
	}

	select {
	case replyMsg := <-replyChan:
		return replyMsg.Msg, true, nil
	case <-time.After(replyTimeout):
		return nil, true, errors.New("Wait for reply timeout")
	}
}
