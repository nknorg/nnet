package chord

import (
	"encoding/hex"

	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/util"
)

// Connect connects to a remote node. optionally with id info to check if
// connection has established
func (c *Chord) Connect(addr string, id []byte) error {
	remoteNode := c.neighbors.GetByID(id)
	if remoteNode != nil {
		log.Infof("Node with id %s is already a neighbor", hex.EncodeToString(id))
		c.addNeighbor(remoteNode)
		return nil
	}
	return c.LocalNode.Connect(addr)
}

// addNeighbor adds a remote node to the neighbor lists of chord overlay
func (c *Chord) addNeighbor(remoteNode *node.RemoteNode) bool {
	if !c.successors.Exists(remoteNode.Id) {
		added, replaced, err := c.successors.AddOrReplace(remoteNode)
		if err != nil {
			log.Error(err)
		}

		if added {
			index := c.successors.GetIndex(remoteNode.Id)
			if index >= 0 {
				for _, f := range c.middlewareStore.successorAdded {
					if !f(remoteNode, index) {
						break
					}
				}
			}
		}

		if replaced != nil {
			for _, f := range c.middlewareStore.successorRemoved {
				if !f(remoteNode) {
					break
				}
			}

			c.maybeRemove(replaced)
		}
	}

	if !c.predecessors.Exists(remoteNode.Id) {
		added, replaced, err := c.predecessors.AddOrReplace(remoteNode)
		if err != nil {
			log.Error(err)
		}

		if added {
			index := c.predecessors.GetIndex(remoteNode.Id)
			if index >= 0 {
				for _, f := range c.middlewareStore.predecessorAdded {
					if !f(remoteNode, index) {
						break
					}
				}
			}
		}

		if replaced != nil {
			for _, f := range c.middlewareStore.predecessorRemoved {
				if !f(remoteNode) {
					break
				}
			}

			c.maybeRemove(replaced)
		}
	}

	for i, finger := range c.fingerTable {
		if !finger.Exists(remoteNode.Id) {
			added, replaced, err := finger.AddOrReplace(remoteNode)
			if err != nil {
				log.Error(err)
			}

			if added {
				index := finger.GetIndex(remoteNode.Id)
				if index >= 0 {
					for _, f := range c.middlewareStore.fingerTableAdded {
						if !f(remoteNode, i, index) {
							break
						}
					}
				}
			}

			if replaced != nil {
				for _, f := range c.middlewareStore.fingerTableRemoved {
					if !f(remoteNode, i) {
						break
					}
				}

				c.maybeRemove(replaced)
			}
		}
	}

	if !c.neighbors.Exists(remoteNode.Id) {
		added, replaced, err := c.neighbors.AddOrReplace(remoteNode)
		if err != nil {
			log.Error(err)
		}

		if added {
			index := c.neighbors.GetIndex(remoteNode.Id)
			if index >= 0 {
				for _, f := range c.middlewareStore.neighborAdded {
					if !f(remoteNode, index) {
						break
					}
				}
			}
		}

		if replaced != nil {
			for _, f := range c.middlewareStore.neighborRemoved {
				if !f(remoteNode) {
					break
				}
			}

			c.maybeRemove(replaced)
		}
	}

	return true
}

// removeNeighbor removes a remote node from the neighbor lists of chord overlay
func (c *Chord) removeNeighbor(remoteNode *node.RemoteNode) bool {
	removed := c.successors.Remove(remoteNode)
	if removed {
		for _, f := range c.middlewareStore.successorRemoved {
			if !f(remoteNode) {
				break
			}
		}
	}

	removed = c.predecessors.Remove(remoteNode)
	if removed {
		for _, f := range c.middlewareStore.predecessorRemoved {
			if !f(remoteNode) {
				break
			}
		}
	}

	for i, finger := range c.fingerTable {
		removed = finger.Remove(remoteNode)
		if removed {
			for _, f := range c.middlewareStore.fingerTableRemoved {
				if !f(remoteNode, i) {
					break
				}
			}
		}
	}

	removed = c.neighbors.Remove(remoteNode)
	if removed {
		for _, f := range c.middlewareStore.neighborRemoved {
			if !f(remoteNode) {
				break
			}
		}
	}

	return true
}

// maybeRemove removes an outbound node that is no longer in successors,
// predecessor, or finger table
func (c *Chord) maybeRemove(remoteNode *node.RemoteNode) bool {
	if !remoteNode.IsOutbound {
		return false
	}

	if c.successors.GetByID(remoteNode.Id) != nil {
		return false
	}

	if c.predecessors.GetByID(remoteNode.Id) != nil {
		return false
	}

	for _, finger := range c.fingerTable {
		if finger.GetByID(remoteNode.Id) != nil {
			return false
		}
	}

	c.removeNeighbor(remoteNode)

	remoteNode.Stop(nil)

	return true
}

func (c *Chord) updateNeighborList(neighborList *NeighborList) error {
	newNodes, err := neighborList.getNewNodesToConnect()
	if err != nil {
		return nil
	}

	errs := util.NewErrors()
	for _, newNode := range newNodes {
		// TODO: check if in connecting list
		err = c.Connect(newNode.Addr, newNode.Id)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs.Merged()
}
