package node

import (
	"errors"
	pbmsg "github.com/nknorg/nnet/protobuf/message"

	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/message"
	"google.golang.org/protobuf/proto"
)

const (
	// msg length is encoded by 32 bit int
	msgLenBytes = 4
)

// RemoteMessage is the received msg from remote node. RemoteNode is nil if
// message is sent by local node.
type RemoteMessage struct {
	RemoteNode *RemoteNode
	Msg        *pbmsg.Message
}

// NewRemoteMessage creates a RemoteMessage with remote node rn and msg
func NewRemoteMessage(rn *RemoteNode, msg *pbmsg.Message) (*RemoteMessage, error) {
	remoteMsg := &RemoteMessage{
		RemoteNode: rn,
		Msg:        msg,
	}
	return remoteMsg, nil
}

// NewPingMessage creates a PING message for heartbeat
func (ln *LocalNode) NewPingMessage() (*pbmsg.Message, error) {
	id, err := message.GenID(ln.MessageIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.Ping{}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_PING,
		RoutingType: pbmsg.RoutingType_DIRECT,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// NewPingReply creates a PING reply for heartbeat
func (ln *LocalNode) NewPingReply(replyToID []byte) (*pbmsg.Message, error) {
	id, err := message.GenID(ln.MessageIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.PingReply{}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_PING,
		RoutingType: pbmsg.RoutingType_DIRECT,
		ReplyToId:   replyToID,
		MessageId:   id,
		Message:     buf,
	}
	return msg, nil
}

// NewExchangeNodeMessage creates a EXCHANGE_NODE message to get node info
func (ln *LocalNode) NewExchangeNodeMessage() (*pbmsg.Message, error) {
	id, err := message.GenID(ln.MessageIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.ExchangeNode{
		Node: ln.Node.Node,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_EXCHANGE_NODE,
		RoutingType: pbmsg.RoutingType_DIRECT,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// NewExchangeNodeReply creates a EXCHANGE_NODE reply to send node info
func (ln *LocalNode) NewExchangeNodeReply(replyToID []byte) (*pbmsg.Message, error) {
	id, err := message.GenID(ln.MessageIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.ExchangeNodeReply{
		Node: ln.Node.Node,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_EXCHANGE_NODE,
		RoutingType: pbmsg.RoutingType_DIRECT,
		ReplyToId:   replyToID,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// NewStopMessage creates a STOP message to notify local node to close
// connection with remote node
func (ln *LocalNode) NewStopMessage() (*pbmsg.Message, error) {
	id, err := message.GenID(ln.MessageIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.Stop{}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_STOP,
		RoutingType: pbmsg.RoutingType_DIRECT,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// handleRemoteMessage handles a remote message and returns error
func (ln *LocalNode) handleRemoteMessage(remoteMsg *RemoteMessage) error {
	if remoteMsg.RemoteNode == nil && remoteMsg.Msg.MessageType != pbmsg.MessageType_BYTES {
		return errors.New("Message is sent by local node")
	}

	switch remoteMsg.Msg.MessageType {
	case pbmsg.MessageType_PING:
		replyMsg, err := ln.NewPingReply(remoteMsg.Msg.MessageId)
		if err != nil {
			return err
		}

		err = remoteMsg.RemoteNode.SendMessageAsync(replyMsg)
		if err != nil {
			return err
		}

	case pbmsg.MessageType_EXCHANGE_NODE:
		msgBody := &pbmsg.ExchangeNode{}
		err := proto.Unmarshal(remoteMsg.Msg.Message, msgBody)
		if err != nil {
			return err
		}

		err = remoteMsg.RemoteNode.setNode(msgBody.Node)
		if err != nil {
			remoteMsg.RemoteNode.Stop(err)
			return err
		}

		replyMsg, err := ln.NewExchangeNodeReply(remoteMsg.Msg.MessageId)
		if err != nil {
			return err
		}

		err = remoteMsg.RemoteNode.SendMessageAsync(replyMsg)
		if err != nil {
			return err
		}

	case pbmsg.MessageType_STOP:
		log.Infof("Received stop message from remote node %v", remoteMsg.RemoteNode)
		remoteMsg.RemoteNode.Stop(nil)

	case pbmsg.MessageType_BYTES:
		msgBody := &pbmsg.Bytes{}
		err := proto.Unmarshal(remoteMsg.Msg.Message, msgBody)
		if err != nil {
			return err
		}

		data := msgBody.Data
		var shouldCallNextMiddleware bool
		for _, mw := range ln.middlewareStore.bytesReceived {
			data, shouldCallNextMiddleware = mw.Func(data, remoteMsg.Msg.MessageId, remoteMsg.Msg.SrcId, remoteMsg.RemoteNode)
			if !shouldCallNextMiddleware {
				break
			}
		}

	default:
		return errors.New("Unknown message type")
	}

	return nil
}
