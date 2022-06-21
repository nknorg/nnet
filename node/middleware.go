package node

import (
	"errors"
	"net"

	"github.com/nknorg/nnet/middleware"
	"github.com/nknorg/nnet/protobuf"
)

// BytesReceived is called when local node receive user-defined BYTES message.
// Message with the same message id will only trigger this middleware once. The
// argument it accepts are bytes data, message ID (can be used to reply
// message), sender ID, and the neighbor that passes you the message (may be
// different from the message sneder). Returns the bytes data to be passed in
// the next middleware and if we should proceed to the next middleware.
type BytesReceived struct {
	Func     func(data, msgID, srcID []byte, remoteNode *RemoteNode) ([]byte, bool)
	Priority int32
}

// LocalNodeWillStart is called right before local node starts listening and
// handling messages. It can be used to add additional data to local node, etc.
// Returns if we should proceed to the next middleware.
type LocalNodeWillStart struct {
	Func     func(*LocalNode) bool
	Priority int32
}

// LocalNodeStarted is called right after local node starts listening and
// handling messages. Returns if we should proceed to the next middleware.
type LocalNodeStarted struct {
	Func     func(*LocalNode) bool
	Priority int32
}

// LocalNodeWillStop is called right before local node stops listening and
// handling messages. Returns if we should proceed to the next middleware.
type LocalNodeWillStop struct {
	Func     func(*LocalNode) bool
	Priority int32
}

// LocalNodeStopped is called right after local node stops listening and
// handling messages. Returns if we should proceed to the next middleware.
type LocalNodeStopped struct {
	Func     func(*LocalNode) bool
	Priority int32
}

// WillConnectToNode is called before local node connect to a new remote node.
// Returns if local node should continue connecting and if we should proceed to
// the next middleware.
type WillConnectToNode struct {
	Func     func(*protobuf.Node) (bool, bool)
	Priority int32
}

// ConnectionAccepted is called when a network connection is accepted. Returns
// if we should accept the connection and if we should proceed to the next
// middleware.
type ConnectionAccepted struct {
	Func     func(net.Conn) (bool, bool)
	Priority int32
}

// RemoteNodeConnected is called when a connection is established with a remote
// node, but the remote node id is typically nil, so it's not a good time to use
// the node yet, but can be used to stop the connection to remote node. Returns
// if we should proceed to the next middleware.
type RemoteNodeConnected struct {
	Func     func(*RemoteNode) bool
	Priority int32
}

// RemoteNodeReady is called when local node has received the node info from
// remote node and the remote node is ready to use. Returns if we should proceed
// to the next middleware.
type RemoteNodeReady struct {
	Func     func(*RemoteNode) bool
	Priority int32
}

// RemoteNodeDisconnected is called when connection to remote node is closed.
// The cause of the connection close can be on either local node or remote node.
// Returns if we should proceed to the next middleware.
type RemoteNodeDisconnected struct {
	Func     func(*RemoteNode) bool
	Priority int32
}

// MessageEncoded is called when a protobuf.Message is encoded into bytes and is
// about to be sent to RemoteNode. This is a good place for transcoding. Returns
// the bytes to send and whether we should proceed to the next middleware. If
// returned bytes is nil, msg will be dropped.
type MessageEncoded struct {
	Func     func(*RemoteNode, []byte) ([]byte, bool)
	Priority int32
}

// MessageWillDecode is called when bytes is received from a RemoteNode and is
// about to be decoded into protobuf.Message. This is a good place for
// transcoding. Returns the bytes to send and whether we should proceed to the
// next middleware. If returned bytes is nil, msg will be dropped.
type MessageWillDecode struct {
	Func     func(*RemoteNode, []byte) ([]byte, bool)
	Priority int32
}

// middlewareStore stores the functions that will be called when certain events
// are triggered or in some pipeline
type middlewareStore struct {
	bytesReceived          []BytesReceived
	localNodeWillStart     []LocalNodeWillStart
	localNodeStarted       []LocalNodeStarted
	localNodeWillStop      []LocalNodeWillStop
	localNodeStopped       []LocalNodeStopped
	willConnectToNode      []WillConnectToNode
	connectionAccepted     []ConnectionAccepted
	remoteNodeConnected    []RemoteNodeConnected
	remoteNodeReady        []RemoteNodeReady
	remoteNodeDisconnected []RemoteNodeDisconnected
	messageEncoded         []MessageEncoded
	messageWillDecode      []MessageWillDecode
}

// newMiddlewareStore creates a middlewareStore
func newMiddlewareStore() *middlewareStore {
	return &middlewareStore{
		bytesReceived:          make([]BytesReceived, 0),
		localNodeWillStart:     make([]LocalNodeWillStart, 0),
		localNodeStarted:       make([]LocalNodeStarted, 0),
		localNodeWillStop:      make([]LocalNodeWillStop, 0),
		localNodeStopped:       make([]LocalNodeStopped, 0),
		willConnectToNode:      make([]WillConnectToNode, 0),
		connectionAccepted:     make([]ConnectionAccepted, 0),
		remoteNodeConnected:    make([]RemoteNodeConnected, 0),
		remoteNodeReady:        make([]RemoteNodeReady, 0),
		remoteNodeDisconnected: make([]RemoteNodeDisconnected, 0),
		messageEncoded:         make([]MessageEncoded, 0),
		messageWillDecode:      make([]MessageWillDecode, 0),
	}
}

// ApplyMiddleware add a middleware to the store
func (store *middlewareStore) ApplyMiddleware(mw interface{}) error {
	switch mw := mw.(type) {
	case BytesReceived:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.bytesReceived = append(store.bytesReceived, mw)
		middleware.Sort(store.bytesReceived)
	case LocalNodeWillStart:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.localNodeWillStart = append(store.localNodeWillStart, mw)
		middleware.Sort(store.localNodeWillStart)
	case LocalNodeStarted:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.localNodeStarted = append(store.localNodeStarted, mw)
		middleware.Sort(store.localNodeStarted)
	case LocalNodeWillStop:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.localNodeWillStop = append(store.localNodeWillStop, mw)
		middleware.Sort(store.localNodeWillStop)
	case LocalNodeStopped:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.localNodeStopped = append(store.localNodeStopped, mw)
		middleware.Sort(store.localNodeStopped)
	case WillConnectToNode:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.willConnectToNode = append(store.willConnectToNode, mw)
		middleware.Sort(store.willConnectToNode)
	case ConnectionAccepted:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.connectionAccepted = append(store.connectionAccepted, mw)
		middleware.Sort(store.connectionAccepted)
	case RemoteNodeConnected:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.remoteNodeConnected = append(store.remoteNodeConnected, mw)
		middleware.Sort(store.remoteNodeConnected)
	case RemoteNodeReady:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.remoteNodeReady = append(store.remoteNodeReady, mw)
		middleware.Sort(store.remoteNodeReady)
	case RemoteNodeDisconnected:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.remoteNodeDisconnected = append(store.remoteNodeDisconnected, mw)
		middleware.Sort(store.remoteNodeDisconnected)
	case MessageEncoded:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.messageEncoded = append(store.messageEncoded, mw)
		middleware.Sort(store.messageEncoded)
	case MessageWillDecode:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.messageWillDecode = append(store.messageWillDecode, mw)
		middleware.Sort(store.messageWillDecode)
	default:
		return errors.New("unknown middleware type")
	}

	return nil
}
