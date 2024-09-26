package overlay

import (
	pbmsg "github.com/nknorg/nnet/protobuf/message"
	"time"

	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/overlay/routing"
)

// Network is the overlay network interface
type Network interface {
	Start(isCreate bool) error
	Stop(error)
	Join(seedNodeAddr string) error
	GetLocalNode() *node.LocalNode
	GetRouters() []routing.Router
	ApplyMiddleware(interface{}) error
	SendMessageAsync(msg *pbmsg.Message, routingType pbmsg.RoutingType) (success bool, err error)
	SendMessageSync(msg *pbmsg.Message, routingType pbmsg.RoutingType, replyTimeout time.Duration) (reply *pbmsg.Message, success bool, err error)
}
