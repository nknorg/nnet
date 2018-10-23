// This example shows how to send and receive messages of different routing types
// Run with default options: go run main.go
// Show usage: go run main.go -h
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
	"github.com/nknorg/nnet/overlay/chord"
	"github.com/nknorg/nnet/overlay/routing"
	"github.com/nknorg/nnet/protobuf"
	"github.com/nknorg/nnet/util"
)

func create(transport string, port uint16, id []byte) (*nnet.NNet, error) {
	conf := config.Config{
		Port:                  port,
		Transport:             transport,
		BaseStabilizeInterval: 233 * time.Millisecond,
	}

	nn, err := nnet.NewNNet(id, conf)
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func main() {
	numNodesPtr := flag.Int("n", 10, "number of nodes")
	transportPtr := flag.String("t", "tcp", "transport type, tcp or kcp")
	flag.Parse()

	if *numNodesPtr < 1 {
		log.Error("Number of nodes must be greater than 0")
		return
	}

	const createPort uint16 = 23333
	var nn *nnet.NNet
	var id []byte
	var err error

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

		err = nn.ApplyMiddleware(routing.RemoteMessageReceived(func(remoteMessage *node.RemoteMessage) (*node.RemoteMessage, bool) {
			if remoteMessage.Msg.MessageType == protobuf.BYTES && remoteMessage.Msg.ReplyToId == nil {
				msgBody := &protobuf.Bytes{}
				err = proto.Unmarshal(remoteMessage.Msg.Message, msgBody)
				if err != nil {
					log.Error(err)
				}

				switch remoteMessage.Msg.RoutingType {
				case protobuf.DIRECT:
					log.Infof("Receive direct message \"%s\" from %x", string(msgBody.Data), remoteMessage.RemoteNode.Id)

				case protobuf.RELAY:
					log.Infof("Receive relay message \"%s\" from %x", string(msgBody.Data), remoteMessage.Msg.SrcId)
					ok, err := nn.SendBytesRelayReply(remoteMessage.Msg.MessageId, []byte("Well received!"), remoteMessage.Msg.SrcId)
					if !ok || err != nil {
						log.Error(err)
					}

				case protobuf.BROADCAST:
					log.Infof("Receive broadcast message \"%s\" from %x", string(msgBody.Data), remoteMessage.Msg.SrcId)
				}

				return nil, false
			}
			return remoteMessage, true
		}))
		if err != nil {
			log.Error(err)
			return
		}

		nnets = append(nnets, nn)
	}

	err = nnets[0].ApplyMiddleware(chord.FingerTableAdded(func(remoteNode *node.RemoteNode, fingerIndex, nodeIndex int) bool {
		err := nnets[0].SendBytesDirectAsync([]byte("Hello my finger!"), remoteNode)
		if err != nil {
			log.Error(err)
		}
		return true
	}))
	if err != nil {
		log.Error(err)
		return
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
	ok, err := nnets[0].SendBytesBroadcastAsync([]byte("This message should be received by EVERYONE multiple times!"))
	if !ok || err != nil {
		log.Error(err)
		return
	}

	time.Sleep(time.Second)
	for i := 3; i > 0; i-- {
		log.Infof("Sending relay message in %d seconds", i)
		time.Sleep(time.Second)
	}
	reply, ok, err := nnets[0].SendBytesRelaySync([]byte("This message should only be received by SOMEONE!"), id)
	if !ok || err != nil {
		log.Error(err)
		return
	}
	msgBody := &protobuf.Bytes{}
	err = proto.Unmarshal(reply.Message, msgBody)
	if err != nil {
		log.Error(err)
	}
	log.Infof("Receive reply message \"%s\" from %x", string(msgBody.Data), reply.SrcId)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Info("\nReceived an interrupt, stopping...\n")

	for i := 0; i < len(nnets); i++ {
		nnets[len(nnets)-1-i].Stop(nil)
	}
}
