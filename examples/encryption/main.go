// This example shows how to encrypt the message between nodes. In this example
// we use a fixed AES key. In production it is recommended to derive the shared
// key per node, or better, per session.

// Run with default options: go run main.go

// Show usage: go run main.go -h

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"flag"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/nknorg/nnet"
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/util"
)

func create(transport string, port uint16) (*nnet.NNet, error) {
	conf := &nnet.Config{
		Port:      port,
		Transport: transport,
	}

	nn, err := nnet.NewNNet(nil, conf)
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func main() {
	transportPtr := flag.String("t", "tcp", "transport type, tcp or kcp")
	numNodesPtr := flag.Int("n", 2, "number of nodes")
	flag.Parse()

	if *numNodesPtr < 1 {
		log.Error("Number of nodes must be greater than 0")
		return
	}

	const createPort uint16 = 23333
	nnets := make([]*nnet.NNet, 0)

	key := sha256.Sum256([]byte{232, 162, 171, 229, 143, 145, 231, 142, 176, 228, 186, 134, 228, 185, 136})
	block, err := aes.NewCipher(key[:])
	if err != nil {
		log.Error(err)
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error(err)
		return
	}

	for i := 0; i < *numNodesPtr; i++ {
		nn, err := create(*transportPtr, createPort+uint16(i))
		if err != nil {
			log.Error(err)
			return
		}

		nnets = append(nnets, nn)

		nn.MustApplyMiddleware(node.MessageEncoded{func(remoteNode *node.RemoteNode, msg []byte) ([]byte, bool) {
			nonce, err := util.RandBytes(aesgcm.NonceSize())
			if err != nil {
				log.Error(err)
				return nil, false
			}
			ciphertext := aesgcm.Seal(nil, nonce, msg, nil)
			msgWithNonce := append(nonce, ciphertext...)
			log.Infof("Encrypted %d bytes", len(msg))
			return msgWithNonce, true
		}, 0})

		nn.MustApplyMiddleware(node.MessageWillDecode{func(remoteNode *node.RemoteNode, msg []byte) ([]byte, bool) {
			if len(msg) < aesgcm.NonceSize() {
				log.Errorf("msg too short")
				return nil, false
			}
			nonce := msg[:aesgcm.NonceSize()]
			plaintext, err := aesgcm.Open(nil, nonce, msg[aesgcm.NonceSize():], nil)
			if err != nil {
				log.Error(err)
				return nil, false
			}
			log.Infof("Decrypted %d bytes", len(plaintext))
			return plaintext, true
		}, 0})
	}

	for i := 0; i < len(nnets); i++ {
		time.Sleep(112358 * time.Microsecond)

		err := nnets[i].Start(i == 0)
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
}
