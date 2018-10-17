package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/nknorg/nnet"
	"github.com/nknorg/nnet/config"
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/overlay/routing"
	"github.com/nknorg/nnet/protobuf"
	"github.com/nknorg/nnet/util"
)

func create(transport string, port uint16, id []byte) (*nnet.NNet, error) {
	conf := config.Config{
		Transport:            transport,
		Port:                 port,
		MinNumSuccessors:     8,
		MinStabilizeInterval: 100 * time.Millisecond,
		MaxStabilizeInterval: 200 * time.Millisecond,
	}

	nn, err := nnet.NewNNet(id, conf)
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
	var id []byte
	var err error

	nnets := make([]*nnet.NNet, 0)

	for i := 0; i < *numNodesPtr; i++ {
		id, err = util.RandBytes(32)
		if err != nil {
			log.Error(err)
			return
		}

		nn, err := create(*transportPtr, createPort+uint16(i), id)
		if err != nil {
			log.Error(err)
			return
		}

		err = nn.ApplyMiddleware(routing.RemoteMessageRouted(func(remoteMessage *node.RemoteMessage, localNode *node.LocalNode, remoteNodes []*node.RemoteNode) (*node.RemoteMessage, *node.LocalNode, []*node.RemoteNode, bool) {
			if remoteMessage.Msg.MessageType == protobuf.BYTES {
				msgBody := &protobuf.Bytes{}
				err = proto.Unmarshal(remoteMessage.Msg.Message, msgBody)
				if err != nil {
					log.Error(err)
				}
				if localNode != nil {
					switch remoteMessage.Msg.RoutingType {
					case protobuf.RELAY:
						log.Infof("Receive relay message \"%s\" from %x", string(msgBody.Data), remoteMessage.Msg.SrcId)
					case protobuf.BROADCAST:
						log.Infof("Receive broadcast message \"%s\" from %x", string(msgBody.Data), remoteMessage.Msg.SrcId)
					}
				}
				// message will not be handle by local node, but will be routed to other remote nodes
				return remoteMessage, nil, remoteNodes, true
			}
			return remoteMessage, localNode, remoteNodes, true
		}))
		if err != nil {
			log.Error(err)
			return
		}

		nnets = append(nnets, nn)
	}

	for i := 0; i < len(nnets); i++ {
		time.Sleep(112358 * time.Microsecond)

		err := nnets[i].Start()
		if err != nil {
			log.Error(err)
			return
		}

		if i > 0 {
			err = nnets[i].Join(fmt.Sprintf("127.0.0.1:%d", createPort))
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
	nnets[0].SendBytesBroadcastAsync([]byte("Hello world!"))

	time.Sleep(time.Second)
	for i := 3; i > 0; i-- {
		log.Infof("Sending relay message in %d seconds", i)
		time.Sleep(time.Second)
	}
	nnets[0].SendBytesRelayAsync([]byte("Hello world!"), id)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Info("\nReceived an interrupt, stopping...\n")

	for i := 0; i < len(nnets); i++ {
		nnets[len(nnets)-1-i].Stop(nil)
	}
}
