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

func create(transport string, port uint16) (*nnet.NNet, error) {
	conf := config.Config{
		Transport:            transport,
		Port:                 port,
		MinNumSuccessors:     8,
		MinStabilizeInterval: 100 * time.Millisecond,
		MaxStabilizeInterval: 200 * time.Millisecond,
	}

	nn, err := nnet.NewNNet(conf)
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func main() {
	numNodesPtr := flag.Int("n", 10, "number of nodes")
	transportPtr := flag.String("t", "tcp", "transport type, e.g. tcp or kcp")
	flag.Parse()

	if *numNodesPtr < 1 {
		log.Error("Number of nodes must be greater than 0")
		return
	}

	const createPort uint16 = 23333
	nnets := make([]*nnet.NNet, 0)

	nn, err := create(*transportPtr, createPort)
	if err != nil {
		log.Error(err)
		return
	}

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

	err = nn.Start()
	if err != nil {
		log.Error(err)
		return
	}

	nnets = append(nnets, nn)

	for i := 0; i < *numNodesPtr-1; i++ {
		time.Sleep(112358 * time.Microsecond)

		nn, err := create(*transportPtr, createPort+uint16(i)+1)
		if err != nil {
			log.Error(err)
			return
		}

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

		err = nn.Start()
		if err != nil {
			log.Error(err)
			return
		}

		err = nn.Join(fmt.Sprintf("127.0.0.1:%d", createPort))
		if err != nil {
			log.Error(err)
			return
		}

		nnets = append(nnets, nn)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Info("\nReceived an interrupt, stopping...\n")

	for i := 0; i < len(nnets); i++ {
		nnets[len(nnets)-1-i].Stop(nil)
	}
}
