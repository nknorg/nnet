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
	localNode *node.LocalNode
	overlay   overlay.Interface
	config    *config.Config
}

// NewNNet creates a new nnet using the userConf provided
func NewNNet(userConf config.Config) (*NNet, error) {
	conf, err := config.MergedConfig(userConf)
	if err != nil {
		return nil, err
	}

	id, err := util.RandBytes(int(conf.NodeIDBytes))
	if err != nil {
		return nil, err
	}

	localNode, err := node.NewLocalNode(id, conf)
	if err != nil {
		return nil, err
	}

	ovl, err := chord.NewChord(localNode, conf)
	if err != nil {
		return nil, err
	}

	nn := &NNet{
		localNode: localNode,
		overlay:   ovl,
		config:    conf,
	}

	return nn, nil
}

// Start starts the lifecycle methods of nnet
func (nn *NNet) Start() error {
	err := nn.localNode.Start()
	if err != nil {
		return err
	}

	err = nn.overlay.Start()
	if err != nil {
		return err
	}

	return nil
}

// Join joins a seed node with address of the form ip:port
func (nn *NNet) Join(seedNodeAddr string) error {
	return nn.overlay.Join(seedNodeAddr)
}
