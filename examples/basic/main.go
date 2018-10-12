package main

import (
	"fmt"
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
	var createPort uint16 = 23333

	_, err := create(createPort)
	if err != nil {
		log.Error(err)
		return
	}

	for i := 0; i < 10; i++ {
		_, err = join(createPort+uint16(i)+1, createPort)
		if err != nil {
			log.Error(err)
		}
	}

	time.Sleep(time.Hour)
}
