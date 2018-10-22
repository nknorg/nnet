package chord

import (
	"errors"
	"sync"

	"github.com/nknorg/nnet/node"
)

// SuccessorAdded is called when a new remote node has been added to the
// successor list. This does not necessarily means the remote node has just
// established a new connection with local node, but may also because another
// remote node previously in successor list was disconnected or removed. When
// being called, it will also pass the index of the remote node in successor
// list after it is added. Returns if we should proceed to the next middleware.
type SuccessorAdded func(*node.RemoteNode, int) bool

// SuccessorRemoved is called when a remote node has been removed from the
// successor list. This does not necessarily means the remote node has just
// disconnected with local node, but may also because another remote node is
// added to the successor list. Returns if we should proceed to the next
// middleware.
type SuccessorRemoved func(*node.RemoteNode) bool

// PredecessorAdded is called when a new remote node has been added to the
// predecessor list. This does not necessarily means the remote node has just
// established a new connection with local node, but may also because another
// remote node previously in predecessor list was disconnected or removed. When
// being called, it will also pass the index of the remote node in predecessor
// list after it is added. Returns if we should proceed to the next middleware.
type PredecessorAdded func(*node.RemoteNode, int) bool

// PredecessorRemoved is called when a remote node has been removed from the
// predecessor list. This does not necessarily means the remote node has just
// disconnected with local node, but may also because another remote node is
// added to the predecessor list. Returns if we should proceed to the next
// middleware.
type PredecessorRemoved func(*node.RemoteNode) bool

// FingerTableAdded is called when a new remote node has been added to the
// finger table. This does not necessarily means the remote node has just
// established a new connection with local node, but may also because another
// remote node previously in finger table was disconnected or removed. When
// being called, it will also pass the index of finger table and the index of
// the remote node in that finger table after it is added. Returns if we should
// proceed to the next middleware.
type FingerTableAdded func(*node.RemoteNode, int, int) bool

// FingerTableRemoved is called when a remote node has been removed from the
// finger table. This does not necessarily means the remote node has just
// disconnected with local node, but may also because another remote node is
// added to the finger table. Returns if we should proceed to the next
// middleware.
type FingerTableRemoved func(*node.RemoteNode, int) bool

// NeighborAdded is called when a new remote node has been added to the neighbor
// list. When being called, it will also pass the index of the remote node in
// neighbor list after it is added. Returns if we should proceed to the next
// middleware.
type NeighborAdded func(*node.RemoteNode, int) bool

// NeighborRemoved is called when a remote node has been removed from the
// neighbor list. Returns if we should proceed to the next middleware.
type NeighborRemoved func(*node.RemoteNode) bool

// middlewareStore stores the functions that will be called when certain events
// are triggered or in some pipeline
type middlewareStore struct {
	sync.RWMutex
	successorAdded     []SuccessorAdded
	successorRemoved   []SuccessorRemoved
	predecessorAdded   []PredecessorAdded
	predecessorRemoved []PredecessorRemoved
	fingerTableAdded   []FingerTableAdded
	fingerTableRemoved []FingerTableRemoved
	neighborAdded      []NeighborAdded
	neighborRemoved    []NeighborRemoved
}

// newMiddlewareStore creates a middlewareStore
func newMiddlewareStore() *middlewareStore {
	return &middlewareStore{
		successorAdded:     make([]SuccessorAdded, 0),
		successorRemoved:   make([]SuccessorRemoved, 0),
		predecessorAdded:   make([]PredecessorAdded, 0),
		predecessorRemoved: make([]PredecessorRemoved, 0),
		fingerTableAdded:   make([]FingerTableAdded, 0),
		fingerTableRemoved: make([]FingerTableRemoved, 0),
		neighborAdded:      make([]NeighborAdded, 0),
		neighborRemoved:    make([]NeighborRemoved, 0),
	}
}

// ApplyMiddleware add a middleware to the store
func (store *middlewareStore) ApplyMiddleware(f interface{}) error {
	store.Lock()
	defer store.Unlock()

	switch f := f.(type) {
	case SuccessorAdded:
		if f == nil {
			return errors.New("middleware is nil")
		}
		store.successorAdded = append(store.successorAdded, f)
	case SuccessorRemoved:
		if f == nil {
			return errors.New("middleware is nil")
		}
		store.successorRemoved = append(store.successorRemoved, f)
	case PredecessorAdded:
		if f == nil {
			return errors.New("middleware is nil")
		}
		store.predecessorAdded = append(store.predecessorAdded, f)
	case PredecessorRemoved:
		if f == nil {
			return errors.New("middleware is nil")
		}
		store.predecessorRemoved = append(store.predecessorRemoved, f)
	case FingerTableAdded:
		if f == nil {
			return errors.New("middleware is nil")
		}
		store.fingerTableAdded = append(store.fingerTableAdded, f)
	case FingerTableRemoved:
		if f == nil {
			return errors.New("middleware is nil")
		}
		store.fingerTableRemoved = append(store.fingerTableRemoved, f)
	case NeighborAdded:
		if f == nil {
			return errors.New("middleware is nil")
		}
		store.neighborAdded = append(store.neighborAdded, f)
	case NeighborRemoved:
		if f == nil {
			return errors.New("middleware is nil")
		}
		store.neighborRemoved = append(store.neighborRemoved, f)
	default:
		return errors.New("unknown middleware type")
	}

	return nil
}
