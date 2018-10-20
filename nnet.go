package nnet

import (
	"github.com/nknorg/nnet/config"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/overlay"
	"github.com/nknorg/nnet/overlay/chord"
	"github.com/nknorg/nnet/util"
)

// NNet is is a peer to peer network
type NNet struct {
	LocalNode *node.LocalNode
	Overlay   overlay.Interface
	Config    *config.Config
}

// NewNNet creates a new nnet using the userConf provided
func NewNNet(id []byte, userConf config.Config) (*NNet, error) {
	conf, err := config.MergedConfig(userConf)
	if err != nil {
		return nil, err
	}

	if id == nil || len(id) == 0 {
		id, err = util.RandBytes(int(conf.NodeIDBytes))
		if err != nil {
			return nil, err
		}
	}

	localNode, err := node.NewLocalNode(id[:], conf)
	if err != nil {
		return nil, err
	}

	ovl, err := chord.NewChord(localNode, conf)
	if err != nil {
		return nil, err
	}

	nn := &NNet{
		LocalNode: localNode,
		Overlay:   ovl,
		Config:    conf,
	}

	return nn, nil
}

// Start starts the lifecycle methods of nnet
func (nn *NNet) Start() error {
	err := nn.LocalNode.Start()
	if err != nil {
		return err
	}

	err = nn.Overlay.Start()
	if err != nil {
		return err
	}

	return nil
}

// Stop stops the lifecycle methods of nnet
func (nn *NNet) Stop(err error) {
	nn.Overlay.Stop(err)
	nn.LocalNode.Stop(err)
}

// Join joins a seed node with address of the form ip:port
func (nn *NNet) Join(seedNodeAddr string) error {
	return nn.Overlay.Join(seedNodeAddr)
}
