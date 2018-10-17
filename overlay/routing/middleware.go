package routing

import (
	"errors"
	"sync"

	"github.com/nknorg/nnet/node"
)

// RemoteMessageReceived is called when a new remote message is received and
// prepare to be handled by the corresponding router. This can be used to
// process, modify or discard message. Returns the remote message to be used (or
// nil to discard the message) and if we should proceed to the next middleware.
type RemoteMessageReceived func(*node.RemoteMessage) (*node.RemoteMessage, bool)

// RemoteMessageRouted is called when the router has computed the node to route
// (could be the local node, remote nodes, or both), and before the message is
// dispatched to local or remote nodes. This can be used to process, modify or
// discard message, or change the computed routes. Returns the remote message to
// be used (or nil to discard the message), local node and remote nodes where
// the message should be routed to, and if we should proceed to the next
// middleware.
type RemoteMessageRouted func(*node.RemoteMessage, *node.LocalNode, []*node.RemoteNode) (*node.RemoteMessage, *node.LocalNode, []*node.RemoteNode, bool)

// middlewareStore stores the functions that will be called when certain events
// are triggered or in some pipeline
type middlewareStore struct {
	sync.RWMutex
	remoteMessageReceived []RemoteMessageReceived
	remoteMessageRouted   []RemoteMessageRouted
}

// newMiddlewareStore creates a middlewareStore
func newMiddlewareStore() *middlewareStore {
	return &middlewareStore{
		remoteMessageReceived: make([]RemoteMessageReceived, 0),
		remoteMessageRouted:   make([]RemoteMessageRouted, 0),
	}
}

// ApplyMiddleware add a middleware to the store
func (store *middlewareStore) ApplyMiddleware(f interface{}) error {
	store.Lock()
	defer store.Unlock()

	switch f := f.(type) {
	case RemoteMessageReceived:
		if f == nil {
			return errors.New("middleware is nil")
		}
		store.remoteMessageReceived = append(store.remoteMessageReceived, f)
	case RemoteMessageRouted:
		if f == nil {
			return errors.New("middleware is nil")
		}
		store.remoteMessageRouted = append(store.remoteMessageRouted, f)
	default:
		return errors.New("unknown middleware type")
	}
	return nil
}
