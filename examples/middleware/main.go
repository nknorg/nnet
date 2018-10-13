package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/nknorg/nnet"
	"github.com/nknorg/nnet/config"
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/overlay/chord"
)

func create(port uint16) (*nnet.NNet, error) {
	conf := config.Config{
		Transport:            "tcp",
		Port:                 port,
		MinNumSuccessors:     8,
		MinStabilizeInterval: 100 * time.Millisecond,
		MaxStabilizeInterval: 200 * time.Millisecond,
	}

	nn, err := nnet.NewNNet(conf)
	if err != nil {
		return nil, err
	}

	nn.Start()

	return nn, nil
}

func join(localPort, seedPort uint16) (*nnet.NNet, error) {
	nn, err := create(localPort)
	if err != nil {
		return nil, err
	}

	err = nn.Join(fmt.Sprintf("127.0.0.1:%d", seedPort))
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func main() {
	numNodesPtr := flag.Int("n", 10, "number of nodes")
	flag.Parse()

	if *numNodesPtr < 1 {
		log.Error("Number of nodes must be greater than 0")
		return
	}

	const createPort uint16 = 23333
	nnets := make([]*nnet.NNet, 0)

	nn, err := create(createPort)
	if err != nil {
		log.Error(err)
		return
	}
	nnets = append(nnets, nn)

	err = nn.ApplyMiddleware(node.RemoteNodeConnected(func(remoteNode *node.RemoteNode) bool {
		log.Infof("Remote node connected: %v", remoteNode)
		return true
	}))
	if err != nil {
		log.Error(err)
		return
	}

	err = nn.ApplyMiddleware(node.RemoteNodeReady(func(remoteNode *node.RemoteNode) bool {
		log.Infof("Remote node ready: %v", remoteNode)
		return true
	}))
	if err != nil {
		log.Error(err)
		return
	}

	err = nn.ApplyMiddleware(node.RemoteNodeDisconnected(func(remoteNode *node.RemoteNode) bool {
		log.Infof("Remote node disconnected: %v", remoteNode)
		return true
	}))
	if err != nil {
		log.Error(err)
		return
	}

	err = nn.ApplyMiddleware(chord.NeighborAdded(func(remoteNode *node.RemoteNode, index int) bool {
		log.Infof("New neighbor %d: %v", index, remoteNode)
		return true
	}))
	if err != nil {
		log.Error(err)
		return
	}

	err = nn.ApplyMiddleware(chord.SuccessorAdded(func(remoteNode *node.RemoteNode, index int) bool {
		log.Infof("New successor %d: %v", index, remoteNode)
		return true
	}))
	if err != nil {
		log.Error(err)
		return
	}

	err = nn.ApplyMiddleware(chord.PredecessorAdded(func(remoteNode *node.RemoteNode, index int) bool {
		log.Infof("New predecessor %d: %v", index, remoteNode)
		return true
	}))
	if err != nil {
		log.Error(err)
		return
	}

	err = nn.ApplyMiddleware(chord.FingerTableAdded(func(remoteNode *node.RemoteNode, fingerIndex, nodeIndex int) bool {
		log.Infof("New finger table %d-%d: %v", fingerIndex, nodeIndex, remoteNode)
		return true
	}))
	if err != nil {
		log.Error(err)
		return
	}

	for i := 0; i < *numNodesPtr-1; i++ {
		nn, err = join(createPort+uint16(i)+1, createPort)
		if err != nil {
			log.Error(err)
			continue
		}
		nnets = append(nnets, nn)

		err = nn.ApplyMiddleware(node.RemoteNodeConnected(func(remoteNode *node.RemoteNode) bool {
			if rand.Float64() < 0.5 {
				remoteNode.Stop(errors.New("YOU ARE UNLUCKY"))
				// stop propagate to the next middleware
				return false
			}
			return true
		}))
		if err != nil {
			log.Error(err)
			return
		}

		err = nn.ApplyMiddleware(node.RemoteNodeConnected(func(remoteNode *node.RemoteNode) bool {
			log.Infof("Only lucky remote node can get here :)")
			return true
		}))
		if err != nil {
			log.Error(err)
			return
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Info("\nReceived an interrupt, stopping...\n")

	for i := 0; i < len(nnets); i++ {
		nnets[len(nnets)-1-i].Stop(nil)
	}
}
