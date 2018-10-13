package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nknorg/nnet"
	"github.com/nknorg/nnet/config"
	"github.com/nknorg/nnet/log"
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

	for i := 0; i < *numNodesPtr-1; i++ {
		nn, err = join(createPort+uint16(i)+1, createPort)
		if err != nil {
			log.Error(err)
			continue
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
