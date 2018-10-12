package overlay

import (
	"errors"
	"fmt"
	"time"

	"github.com/nknorg/nnet/common"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/overlay/routing"
	"github.com/nknorg/nnet/protobuf"
	"github.com/nknorg/nnet/util"
)

const (
	// Max number of msg to be processed by local node that can be buffered
	localMsgChanLen = 23333

	// Timeout for reply message
	replyTimeout = 5 * time.Second
)

// Interface is the overlay network interface
type Interface interface {
	Start() error
	Stop(error)
	Join(string) error
}

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
		LocalMsgChan: make(chan *node.RemoteMessage, localMsgChanLen),
		routers:      make(map[protobuf.RoutingType]routing.Router),
	}
	return overlay, nil
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

// SendMessage sends msg to the best next hop, returns reply chans if hasReply =
// true and aggregated errors during message sending
func (ovl *Overlay) SendMessage(msg *protobuf.Message, routingType protobuf.RoutingType, hasReply bool) ([]chan *node.RemoteMessage, bool, error) {
	router, err := ovl.GetRouter(routingType)
	if err != nil {
		return nil, false, err
	}

	_, remoteNodes, err := router.GetNodeToRoute(&node.RemoteMessage{Msg: msg})
	if err != nil {
		return nil, false, err
	}

	if remoteNodes == nil || len(remoteNodes) == 0 {
		return nil, false, errors.New("No remote node to route")
	}

	var replyChans []chan *node.RemoteMessage
	errs := util.NewErrors()
	success := false
	if hasReply {
		replyChans = make([]chan *node.RemoteMessage, 0)
		for _, remoteNode := range remoteNodes {
			replyChan, err := remoteNode.SendMessage(msg, true)
			if err != nil {
				errs = append(errs, err)
			} else {
				replyChans = append(replyChans, replyChan)
				success = true
			}
		}
	} else {
		for _, remoteNode := range remoteNodes {
			_, err := remoteNode.SendMessage(msg, false)
			if err != nil {
				errs = append(errs, err)
			} else {
				success = true
			}
		}
	}

	if !success {
		replyChans = nil
	}

	return replyChans, true, errs.Merged()
}

// SendMessageAsync sends msg to the best next hop and returns aggretated error
// during message sending
func (ovl *Overlay) SendMessageAsync(msg *protobuf.Message, routingType protobuf.RoutingType) (bool, error) {
	_, success, err := ovl.SendMessage(msg, routingType, false)
	return success, err
}

// SendMessageSync sends msg to the best next hop, returns reply message and
// aggregated error during message sending, will also returns error if haven't
// receive reply before timeout
func (ovl *Overlay) SendMessageSync(msg *protobuf.Message, routingType protobuf.RoutingType) (*protobuf.Message, bool, error) {
	replyChans, success, err := ovl.SendMessage(msg, routingType, true)
	if !success {
		return nil, success, err
	}

	raceReplyChan := make(chan *node.RemoteMessage)

	for i := range replyChans {
		go func(replyChan chan *node.RemoteMessage) {
			select {
			case replyMsg := <-replyChan:
				select {
				case raceReplyChan <- replyMsg:
				default:
				}
			case <-time.After(replyTimeout):
			}
		}(replyChans[i])
	}

	select {
	case replyMsg := <-raceReplyChan:
		return replyMsg.Msg, true, nil
	case <-time.After(replyTimeout):
		return nil, true, errors.New("Wait for reply timeout")
	}
}
