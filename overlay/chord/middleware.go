package chord

import (
	"errors"

	"github.com/nknorg/nnet/middleware"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/overlay"
)

// SuccessorAdded is called when a new remote node has been added to the
// successor list. This does not necessarily means the remote node has just
// established a new connection with local node, but may also because another
// remote node previously in successor list was disconnected or removed. When
// being called, it will also pass the index of the remote node in successor
// list after it is added. Returns if we should proceed to the next middleware.
type SuccessorAdded struct {
	Func     func(*node.RemoteNode, int) bool
	Priority int32
}

// SuccessorRemoved is called when a remote node has been removed from the
// successor list. This does not necessarily means the remote node has just
// disconnected with local node, but may also because another remote node is
// added to the successor list. Returns if we should proceed to the next
// middleware.
type SuccessorRemoved struct {
	Func     func(*node.RemoteNode) bool
	Priority int32
}

// PredecessorAdded is called when a new remote node has been added to the
// predecessor list. This does not necessarily means the remote node has just
// established a new connection with local node, but may also because another
// remote node previously in predecessor list was disconnected or removed. When
// being called, it will also pass the index of the remote node in predecessor
// list after it is added. Returns if we should proceed to the next middleware.
type PredecessorAdded struct {
	Func     func(*node.RemoteNode, int) bool
	Priority int32
}

// PredecessorRemoved is called when a remote node has been removed from the
// predecessor list. This does not necessarily means the remote node has just
// disconnected with local node, but may also because another remote node is
// added to the predecessor list. Returns if we should proceed to the next
// middleware.
type PredecessorRemoved struct {
	Func     func(*node.RemoteNode) bool
	Priority int32
}

// FingerTableAdded is called when a new remote node has been added to the
// finger table. This does not necessarily means the remote node has just
// established a new connection with local node, but may also because another
// remote node previously in finger table was disconnected or removed. When
// being called, it will also pass the index of finger table and the index of
// the remote node in that finger table after it is added. Returns if we should
// proceed to the next middleware.
type FingerTableAdded struct {
	Func     func(*node.RemoteNode, int, int) bool
	Priority int32
}

// FingerTableRemoved is called when a remote node has been removed from the
// finger table. This does not necessarily means the remote node has just
// disconnected with local node, but may also because another remote node is
// added to the finger table. Returns if we should proceed to the next
// middleware.
type FingerTableRemoved struct {
	Func     func(*node.RemoteNode, int) bool
	Priority int32
}

// NeighborAdded is called when a new remote node has been added to the neighbor
// list. When being called, it will also pass the index of the remote node in
// neighbor list after it is added. Returns if we should proceed to the next
// middleware.
type NeighborAdded struct {
	Func     func(*node.RemoteNode, int) bool
	Priority int32
}

// NeighborRemoved is called when a remote node has been removed from the
// neighbor list. Returns if we should proceed to the next middleware.
type NeighborRemoved struct {
	Func     func(*node.RemoteNode) bool
	Priority int32
}

// middlewareStore stores the functions that will be called when certain events
// are triggered or in some pipeline
type middlewareStore struct {
	networkWillStart   []overlay.NetworkWillStart
	networkStarted     []overlay.NetworkStarted
	networkWillStop    []overlay.NetworkWillStop
	networkStopped     []overlay.NetworkStopped
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
		networkWillStart:   make([]overlay.NetworkWillStart, 0),
		networkStarted:     make([]overlay.NetworkStarted, 0),
		networkWillStop:    make([]overlay.NetworkWillStop, 0),
		networkStopped:     make([]overlay.NetworkStopped, 0),
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
func (store *middlewareStore) ApplyMiddleware(mw interface{}) error {
	switch mw := mw.(type) {
	case overlay.NetworkWillStart:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.networkWillStart = append(store.networkWillStart, mw)
		middleware.Sort(store.networkWillStart)
	case overlay.NetworkStarted:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.networkStarted = append(store.networkStarted, mw)
		middleware.Sort(store.networkStarted)
	case overlay.NetworkWillStop:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.networkWillStop = append(store.networkWillStop, mw)
		middleware.Sort(store.networkWillStop)
	case overlay.NetworkStopped:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.networkStopped = append(store.networkStopped, mw)
		middleware.Sort(store.networkStopped)
	case SuccessorAdded:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.successorAdded = append(store.successorAdded, mw)
		middleware.Sort(store.successorAdded)
	case SuccessorRemoved:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.successorRemoved = append(store.successorRemoved, mw)
		middleware.Sort(store.successorRemoved)
	case PredecessorAdded:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.predecessorAdded = append(store.predecessorAdded, mw)
		middleware.Sort(store.predecessorAdded)
	case PredecessorRemoved:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.predecessorRemoved = append(store.predecessorRemoved, mw)
		middleware.Sort(store.predecessorRemoved)
	case FingerTableAdded:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.fingerTableAdded = append(store.fingerTableAdded, mw)
		middleware.Sort(store.fingerTableAdded)
	case FingerTableRemoved:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.fingerTableRemoved = append(store.fingerTableRemoved, mw)
		middleware.Sort(store.fingerTableRemoved)
	case NeighborAdded:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.neighborAdded = append(store.neighborAdded, mw)
		middleware.Sort(store.neighborAdded)
	case NeighborRemoved:
		if mw.Func == nil {
			return errors.New("middleware function is nil")
		}
		store.neighborRemoved = append(store.neighborRemoved, mw)
		middleware.Sort(store.neighborRemoved)
	default:
		return errors.New("unknown middleware type")
	}

	return nil
}
