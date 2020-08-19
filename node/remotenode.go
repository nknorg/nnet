package node

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/nknorg/nnet/cache"
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/multiplexer"
	"github.com/nknorg/nnet/protobuf"
	"github.com/nknorg/nnet/transport"
	"github.com/nknorg/nnet/util"
)

const (
	// A grace period that allows remote node to send messages in queue
	stopGracePeriod = 100 * time.Millisecond

	// Number of retries to get remote node when remote node starts
	startRetries = 3
)

// RemoteNode is a remote node
type RemoteNode struct {
	*Node
	LocalNode  *LocalNode
	IsOutbound bool
	conn       net.Conn
	rxMsgChan  chan *protobuf.Message
	txMsgChan  chan *protobuf.Message
	txMsgCache cache.Cache

	sync.RWMutex
	lastRxTime    time.Time
	roundTripTime time.Duration
}

// NewRemoteNode creates a remote node
func NewRemoteNode(localNode *LocalNode, conn net.Conn, isOutbound bool, n *protobuf.Node) (*RemoteNode, error) {
	if localNode == nil {
		return nil, errors.New("Local node is nil")
	}
	if conn == nil {
		return nil, errors.New("conn is nil")
	}

	var node *Node
	var err error
	if n != nil {
		node, err = newNode(n)
		if err != nil {
			return nil, err
		}
	} else {
		node, err = NewNode(nil, "")
		if err != nil {
			return nil, err
		}
	}

	txMsgCache := cache.NewGoCache(localNode.RemoteTxMsgCacheExpiration, localNode.RemoteTxMsgCacheCleanupInterval)

	remoteNode := &RemoteNode{
		Node:       node,
		LocalNode:  localNode,
		conn:       conn,
		IsOutbound: isOutbound,
		rxMsgChan:  make(chan *protobuf.Message, localNode.RemoteRxMsgChanLen),
		txMsgChan:  make(chan *protobuf.Message, localNode.RemoteTxMsgChanLen),
		txMsgCache: txMsgCache,
		lastRxTime: time.Now(),
	}

	return remoteNode, nil
}

func (rn *RemoteNode) String() string {
	if !rn.IsReady() {
		return fmt.Sprintf("<%s>", rn.conn.RemoteAddr().String())
	}
	return fmt.Sprintf("%v<%s>", rn.Node, rn.conn.RemoteAddr().String())
}

// GetConn returns the connection with remote node
func (rn *RemoteNode) GetConn() net.Conn {
	return rn.conn
}

// GetRoundTripTime returns the measured round trip time between local node and
// remote node. Will return 0 if no result available yet.
func (rn *RemoteNode) GetRoundTripTime() time.Duration {
	rn.RLock()
	defer rn.RUnlock()
	return rn.roundTripTime
}

// setRoundTripTime sets the measured round trip time between local node and
// remote node.
func (rn *RemoteNode) setRoundTripTime(roundTripTime time.Duration) {
	rn.Lock()
	rn.roundTripTime = roundTripTime
	rn.Unlock()
}

// GetLastRxTime returns the time received last ping response
func (rn *RemoteNode) GetLastRxTime() time.Time {
	rn.RLock()
	defer rn.RUnlock()
	return rn.lastRxTime
}

// setLastRxTime sets the time received last ping response
func (rn *RemoteNode) setLastRxTime(lastRxTime time.Time) {
	rn.Lock()
	rn.lastRxTime = lastRxTime
	rn.Unlock()
}

func (rn *RemoteNode) setNode(n *protobuf.Node) error {
	rn.Node.Lock()
	defer rn.Node.Unlock()

	if rn.Id != nil && !bytes.Equal(rn.Id, n.Id) {
		return fmt.Errorf("Node id %x is different from expected value %x", n.Id, rn.Id)
	}

	remoteAddr, err := transport.Parse(n.Addr, rn.LocalNode.SupportedTransports)
	if err != nil {
		return fmt.Errorf("Parse node addr %s error: %s", n.Addr, err)
	}

	if remoteAddr.Host == "" {
		connAddr := rn.conn.RemoteAddr().String()
		remoteAddr.Host, _, err = net.SplitHostPort(connAddr)
		if err != nil {
			return fmt.Errorf("Parse conn remote addr %s error: %s", connAddr, err)
		}
		n.Addr = remoteAddr.String()
	}

	if rn.Addr != "" {
		expectedAddr, err := transport.Parse(rn.Addr, rn.LocalNode.SupportedTransports)
		if err == nil && expectedAddr.Host == "" {
			connAddr := rn.conn.RemoteAddr().String()
			expectedAddr.Host, _, err = net.SplitHostPort(connAddr)
			if err == nil {
				rn.Addr = expectedAddr.String()
			}
		}

		if rn.Addr != n.Addr {
			return fmt.Errorf("Node addr %s is different from expected value %s", n.Addr, rn.Addr)
		}
	}

	if !proto.Equal(rn.Node.Node, n) {
		rn.Node.Node = n
	}

	return nil
}

// Start starts the runtime loop of the remote node
func (rn *RemoteNode) Start() error {
	rn.StartOnce.Do(func() {
		if rn.IsStopped() {
			return
		}

		go rn.handleMsg()
		go rn.startMultiplexer()
		go rn.startMeasuringRoundTripTime()

		go func() {
			var n *protobuf.Node
			var err error

			for i := 0; i < startRetries; i++ {
				n, err = rn.ExchangeNode()
				if err == nil {
					break
				}
			}
			if err != nil {
				rn.Stop(fmt.Errorf("Get node error: %s", err))
				return
			}

			err = rn.setNode(n)
			if err != nil {
				rn.Stop(err)
				return
			}

			var existing *RemoteNode
			rn.LocalNode.neighbors.Range(func(key, value interface{}) bool {
				remoteNode, ok := value.(*RemoteNode)
				if ok && remoteNode.IsReady() && bytes.Equal(remoteNode.Id, n.Id) {
					if remoteNode.IsStopped() {
						log.Warningf("Remove stopped remote node %v from list", remoteNode)
						rn.LocalNode.neighbors.Delete(key)
					} else {
						existing = remoteNode
					}
					return false
				}
				return true
			})
			if existing != nil {
				rn.Stop(fmt.Errorf("Node with id %x is already connected at addr %s", existing.Id, existing.conn.RemoteAddr().String()))
				return
			}

			rn.SetReady(true)

			for _, mw := range rn.LocalNode.middlewareStore.remoteNodeReady {
				if !mw.Func(rn) {
					break
				}
			}
		}()
	})

	return nil
}

// Stop stops the runtime loop of the remote node
func (rn *RemoteNode) Stop(err error) {
	rn.StopOnce.Do(func() {
		if err != nil {
			log.Warningf("Remote node %v stops because of error: %s", rn, err)
		} else {
			log.Infof("Remote node %v stops", rn)
		}

		err = rn.NotifyStop()
		if err != nil {
			log.Warning("Notify remote node %v stop error:", rn, err)
		}

		time.AfterFunc(stopGracePeriod, func() {
			rn.LifeCycle.Stop()

			if rn.conn != nil {
				rn.LocalNode.neighbors.Delete(rn.conn.RemoteAddr().String())
				rn.conn.Close()
			}

			for _, mw := range rn.LocalNode.middlewareStore.remoteNodeDisconnected {
				if !mw.Func(rn) {
					break
				}
			}
		})
	})
}

func (rn *RemoteNode) startMultiplexer() {
	mux, err := multiplexer.NewMultiplexer(rn.LocalNode.Multiplexer, rn.conn, rn.IsOutbound)
	if err != nil {
		rn.Stop(fmt.Errorf("Create multiplexer error: %s", err))
		return
	}

	var conn net.Conn
	if rn.IsOutbound {
		for i := uint32(0); i < rn.LocalNode.NumStreamsToOpen; i++ {
			conn, err = mux.OpenStream()
			if err != nil {
				rn.Stop(fmt.Errorf("Open stream error: %s", err))
				return
			}
			go rn.rx(conn, i == 0)
		}
	} else {
		for i := uint32(0); i < rn.LocalNode.NumStreamsToAccept; i++ {
			conn, err = mux.AcceptStream()
			if err != nil {
				rn.Stop(fmt.Errorf("Accept stream error: %s", err))
				return
			}
			go rn.rx(conn, true)
		}
	}
}

// handleMsg starts a loop that handles received msg
func (rn *RemoteNode) handleMsg() {
	var msg *protobuf.Message
	var remoteMsg *RemoteMessage
	var msgChan chan *RemoteMessage
	var added, ok bool
	var err error
	keepAliveTimeoutTimer := time.NewTimer(rn.LocalNode.KeepAliveTimeout)

NEXT_MESSAGE:
	for {
		if rn.IsStopped() {
			util.StopTimer(keepAliveTimeoutTimer)
			return
		}

		select {
		case msg, ok = <-rn.rxMsgChan:
			if !ok {
				util.StopTimer(keepAliveTimeoutTimer)
				return
			}

			for _, routingType := range rn.LocalNode.LocalRxMsgCacheRoutingType {
				if routingType == msg.RoutingType {
					added, err = rn.LocalNode.AddToRxCache(msg.MessageId)
					if err != nil {
						log.Errorf("Add msg id %x to rx cache error: %v", msg.MessageId, err)
						continue NEXT_MESSAGE
					}
					if !added {
						continue NEXT_MESSAGE
					}
				}
			}

			remoteMsg, err = NewRemoteMessage(rn, msg)
			if err != nil {
				log.Errorf("New remote message error: %v", err)
				continue
			}

			msgChan, err = rn.LocalNode.GetRxMsgChan(msg.RoutingType)
			if err != nil {
				log.Errorf("Get rx msg chan for routing type %v error: %v", msg.RoutingType, err)
				continue
			}

			select {
			case msgChan <- remoteMsg:
			default:
				log.Warningf("Msg chan full for routing type %d, discarding msg", msg.RoutingType)
			}
		case <-keepAliveTimeoutTimer.C:
			if time.Since(rn.GetLastRxTime()) > rn.LocalNode.KeepAliveTimeout {
				rn.Stop(errors.New("keepalive timeout"))
			}
		}

		util.ResetTimer(keepAliveTimeoutTimer, rn.LocalNode.KeepAliveTimeout)
	}
}

// handleMsgBuf unmarshal buf to msg and send it to msg chan of the local node
func (rn *RemoteNode) handleMsgBuf(buf []byte) {
	msg := &protobuf.Message{}
	err := proto.Unmarshal(buf, msg)
	if err != nil {
		rn.Stop(fmt.Errorf("unmarshal msg error: %s", err))
		return
	}

	select {
	case rn.rxMsgChan <- msg:
	default:
		log.Warning("Rx msg chan full, discarding msg")
	}
}

// rx receives and handle data from RemoteNode rn
func (rn *RemoteNode) rx(conn net.Conn, isActive bool) {
	msgLenBuf := make([]byte, msgLenBytes)
	var readLen uint32

	if isActive {
		go rn.tx(conn)
	}

	for {
		if rn.IsStopped() {
			return
		}

		l, err := conn.Read(msgLenBuf)
		if err != nil {
			rn.Stop(fmt.Errorf("Read msg len error: %s", err))
			continue
		}
		if l != msgLenBytes {
			rn.Stop(fmt.Errorf("Msg len has %d bytes, which is less than expected %d", l, msgLenBytes))
			continue
		}

		if !isActive {
			isActive = true
			go rn.tx(conn)
		}

		msgLen := binary.BigEndian.Uint32(msgLenBuf)
		if msgLen < 0 {
			rn.Stop(fmt.Errorf("Msg len %d overflow", msgLen))
			continue
		}

		if msgLen > rn.LocalNode.MaxMessageSize {
			rn.Stop(fmt.Errorf("Msg size %d exceeds max msg size %d", msgLen, rn.LocalNode.MaxMessageSize))
			continue
		}

		buf := make([]byte, msgLen)

		for readLen = 0; readLen < msgLen; readLen += uint32(l) {
			l, err = conn.Read(buf[readLen:])
			if err != nil {
				break
			}
		}

		if err != nil {
			rn.Stop(fmt.Errorf("Read msg error: %s", err))
			continue
		}

		if readLen > msgLen {
			rn.Stop(fmt.Errorf("Msg has %d bytes, which is more than expected %d", readLen, msgLen))
			continue
		}

		var shouldCallNextMiddleware bool
		for _, mw := range rn.LocalNode.middlewareStore.messageWillDecode {
			buf, shouldCallNextMiddleware = mw.Func(rn, buf)
			if buf == nil || !shouldCallNextMiddleware {
				break
			}
		}

		if buf == nil {
			continue
		}

		rn.handleMsgBuf(buf)
	}
}

// tx marshals and sends data in txMsgChan to RemoteNode rn
func (rn *RemoteNode) tx(conn net.Conn) {
	var msg *protobuf.Message
	var buf []byte
	var ok bool
	var err error
	msgLenBuf := make([]byte, msgLenBytes)
	txTimeoutTimer := time.NewTimer(time.Second)

	for {
		if rn.IsStopped() {
			util.StopTimer(txTimeoutTimer)
			return
		}

		select {
		case msg, ok = <-rn.txMsgChan:
			if !ok {
				util.StopTimer(txTimeoutTimer)
				return
			}

			buf, err = proto.Marshal(msg)
			if err != nil {
				log.Errorf("Marshal msg error: %v", err)
				continue
			}

			var shouldCallNextMiddleware bool
			for _, mw := range rn.LocalNode.middlewareStore.messageEncoded {
				buf, shouldCallNextMiddleware = mw.Func(rn, buf)
				if buf == nil || !shouldCallNextMiddleware {
					break
				}
			}

			if buf == nil {
				continue
			}

			if uint32(len(buf)) > rn.LocalNode.MaxMessageSize {
				log.Errorf("Msg size %d exceeds max msg size %d", len(buf), rn.LocalNode.MaxMessageSize)
				continue
			}

			binary.BigEndian.PutUint32(msgLenBuf, uint32(len(buf)))

			_, err = conn.Write(msgLenBuf)
			if err != nil {
				rn.Stop(fmt.Errorf("Write to conn error: %s", err))
				continue
			}

			_, err = conn.Write(buf)
			if err != nil {
				rn.Stop(fmt.Errorf("Write to conn error: %s", err))
				continue
			}
		case <-txTimeoutTimer.C:
		}

		util.ResetTimer(txTimeoutTimer, time.Second)
	}
}

// startMeasuringRoundTripTime starts to periodically send ping message to
// measure round trip time to remote node.
func (rn *RemoteNode) startMeasuringRoundTripTime() {
	var err error
	var txTime, rxTime time.Time
	var lastRoundTripTime, interval time.Duration

	for {
		time.Sleep(interval)

		if rn.IsStopped() {
			return
		}

		txTime = time.Now()
		err = rn.Ping()
		if err != nil {
			log.Warningf("Ping %v error: %v", rn, err)
			// This will guarantee rn.roundTripTime immediately becomes larger than
			// any other available neighbors
			lastRoundTripTime = rn.LocalNode.DefaultReplyTimeout * 2
			interval = rn.LocalNode.MeasureRoundTripTimeInterval / 5
		} else {
			rxTime = time.Now()
			lastRoundTripTime = rxTime.Sub(txTime)
			rn.setLastRxTime(rxTime)
			interval = util.RandDuration(rn.LocalNode.MeasureRoundTripTimeInterval, 1.0/5.0)
		}

		if rn.roundTripTime > 0 {
			rn.setRoundTripTime((rn.roundTripTime + lastRoundTripTime) / 2)
		} else {
			rn.setRoundTripTime(lastRoundTripTime)
		}
	}
}

// SendMessage marshals and sends msg, will returns a RemoteMessage chan if
// hasReply is true and reply is received within replyTimeout.
func (rn *RemoteNode) SendMessage(msg *protobuf.Message, hasReply bool, replyTimeout time.Duration) (<-chan *RemoteMessage, error) {
	if rn.IsStopped() {
		return nil, errors.New("Remote node has stopped")
	}

	if len(msg.MessageId) == 0 {
		return nil, errors.New("Message ID is empty")
	}

	for _, routingType := range rn.LocalNode.RemoteTxMsgCacheRoutingType {
		if routingType == msg.RoutingType {
			_, found := rn.txMsgCache.Get(msg.MessageId)
			if found {
				return nil, nil
			}

			err := rn.txMsgCache.Add(msg.MessageId, struct{}{})
			if err != nil {
				return nil, err
			}
		}
	}

	select {
	case rn.txMsgChan <- msg:
	default:
		return nil, errors.New("Tx msg chan full, discarding msg")
	}

	if hasReply {
		return rn.LocalNode.AllocReplyChan(msg.MessageId, replyTimeout)
	}

	return nil, nil
}

// SendMessageAsync sends msg and returns if there is an error
func (rn *RemoteNode) SendMessageAsync(msg *protobuf.Message) error {
	_, err := rn.SendMessage(msg, false, 0)
	return err
}

// SendMessageSync sends msg, returns reply message or error if don't receive
// reply within replyTimeout. Will use default reply timeout in config if
// replyTimeout = 0.
func (rn *RemoteNode) SendMessageSync(msg *protobuf.Message, replyTimeout time.Duration) (*RemoteMessage, error) {
	if replyTimeout == 0 {
		replyTimeout = rn.LocalNode.DefaultReplyTimeout
	}

	replyChan, err := rn.SendMessage(msg, true, replyTimeout)
	if err != nil {
		return nil, err
	}

	select {
	case replyMsg := <-replyChan:
		return replyMsg, nil
	case <-time.After(replyTimeout):
		return nil, errors.New("Wait for reply timeout")
	}
}

// Ping sends a Ping message to remote node and wait for reply
func (rn *RemoteNode) Ping() error {
	msg, err := rn.LocalNode.NewPingMessage()
	if err != nil {
		return err
	}

	_, err = rn.SendMessageSync(msg, 0)
	if err != nil {
		return err
	}

	return nil
}

// ExchangeNode sends a ExchangeNode message to remote node and wait for reply
func (rn *RemoteNode) ExchangeNode() (*protobuf.Node, error) {
	msg, err := rn.LocalNode.NewExchangeNodeMessage()
	if err != nil {
		return nil, err
	}

	reply, err := rn.SendMessageSync(msg, 0)
	if err != nil {
		return nil, err
	}

	replyBody := &protobuf.ExchangeNodeReply{}
	err = proto.Unmarshal(reply.Msg.Message, replyBody)
	if err != nil {
		return nil, err
	}

	return replyBody.Node, nil
}

// NotifyStop sends a Stop message to remote node to notify it that we will
// close connection with it
func (rn *RemoteNode) NotifyStop() error {
	msg, err := rn.LocalNode.NewStopMessage()
	if err != nil {
		return err
	}

	err = rn.SendMessageAsync(msg)
	if err != nil {
		return err
	}

	return nil
}
