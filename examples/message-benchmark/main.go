// This example is a simple benchmark measuring the throughput of broadcast
// messages.

// Run with default options: go run main.go

// Show usage: go run main.go -h

package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/nknorg/nnet"
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/protobuf"
	"github.com/nknorg/nnet/util"
)

func create(transport string, port uint16, id []byte) (*nnet.NNet, error) {
	conf := &nnet.Config{
		Port:                port,
		Transport:           transport,
		NumFingerSuccessors: 1,
	}

	nn, err := nnet.NewNNet(id, conf)
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func main() {
	transportPtr := flag.String("t", "tcp", "transport type, tcp or kcp")
	numNodesPtr := flag.Int("n", 2, "number of nodes")
	broadcastTypePtr := flag.String("b", "tree", "broadcast type, push or tree")
	msgSizePtr := flag.Int("m", 1024, "message size in bytes")
	flag.Parse()

	if *numNodesPtr < 2 {
		log.Error("Number of nodes must be greater than 1")
		return
	}

	var broadcastType protobuf.RoutingType
	switch *broadcastTypePtr {
	case "push":
		broadcastType = protobuf.BROADCAST_PUSH
	case "tree":
		broadcastType = protobuf.BROADCAST_TREE
	default:
		log.Error("Unknown broadcast type")
		return
	}

	if *msgSizePtr < 1 {
		log.Error("Msg size must be greater than 0")
		return
	}

	const createPort uint16 = 23333
	var nn *nnet.NNet
	var id []byte
	var err error
	var msgCount int
	var msgCountLock sync.RWMutex

	nnets := make([]*nnet.NNet, 0)

	for i := 0; i < *numNodesPtr; i++ {
		id, err = util.RandBytes(32)
		if err != nil {
			log.Error(err)
			return
		}

		nn, err = create(*transportPtr, createPort+uint16(i), id)
		if err != nil {
			log.Error(err)
			return
		}

		nn.MustApplyMiddleware(node.BytesReceived(func(msg, msgID, srcID []byte, remoteNode *node.RemoteNode) ([]byte, bool) {
			msgCountLock.Lock()
			msgCount++
			msgCountLock.Unlock()
			return msg, true
		}))

		nnets = append(nnets, nn)
	}

	for i := 0; i < len(nnets); i++ {
		time.Sleep(112358 * time.Microsecond)

		err = nnets[i].Start(i == 0)
		if err != nil {
			log.Error(err)
			return
		}

		if i > 0 {
			err = nnets[i].Join(nnets[0].GetLocalNode().Addr)
			if err != nil {
				log.Error(err)
				return
			}
		}
	}

	time.Sleep(time.Duration(*numNodesPtr/5) * time.Second)
	for i := 3; i > 0; i-- {
		log.Infof("Sending broadcast message in %d seconds", i)
		time.Sleep(time.Second)
	}
	log.Info("Starting...")

	msg, _ := util.RandBytes(*msgSizePtr)
	go func() {
		for {
			nnets[0].SendBytesBroadcastAsync(msg, broadcastType)
		}
	}()

	time.Sleep(time.Second)

	msgCountLock.RLock()
	msgCountHistory := []int{msgCount}
	msgCountLock.RUnlock()

	go func() {
		for {
			time.Sleep(time.Second)

			msgCountLock.RLock()
			msgCountHistory = append(msgCountHistory, msgCount)
			msgCountLock.RUnlock()

			msgPerNode := (msgCountHistory[len(msgCountHistory)-1] - msgCountHistory[len(msgCountHistory)-2]) / (len(nnets) - 1)
			log.Infof("Each node receives %d msg/s or %f MB/s", msgPerNode, float32(msgPerNode*len(msg))/1024/1024)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Info("\nReceived an interrupt, stopping...\n")

	var wg sync.WaitGroup
	for i := 0; i < len(nnets); i++ {
		wg.Add(1)
		go func(nn *nnet.NNet) {
			nn.Stop(nil)
			wg.Done()
		}(nnets[len(nnets)-1-i])
	}
	wg.Wait()

	if len(msgCountHistory) > 2 {
		msgPerNode := (msgCountHistory[len(msgCountHistory)-1] - msgCountHistory[0]) / (len(msgCountHistory) - 1) / (len(nnets) - 1)
		log.Infof("Each node receives %d msg/s or %f MB/s on average", msgPerNode, float32(msgPerNode*len(msg))/1024/1024)
	}
}
