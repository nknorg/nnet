package nnet

import (
	"time"

	"github.com/nknorg/nnet/message"
	"github.com/nknorg/nnet/node"
	pbmsg "github.com/nknorg/nnet/protobuf/message"
	"google.golang.org/protobuf/proto"
)

// NewDirectBytesMessage creates a BYTES message that send arbitrary bytes to a
// given remote node
func (nn *NNet) NewDirectBytesMessage(data []byte) (*pbmsg.Message, error) {
	id, err := message.GenID(nn.GetLocalNode().MessageIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.Bytes{
		Data: data,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_BYTES,
		RoutingType: pbmsg.RoutingType_DIRECT,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// NewRelayBytesMessage creates a BYTES message that send arbitrary bytes to the
// remote node that has the smallest distance to a given key
func (nn *NNet) NewRelayBytesMessage(data, srcID, key []byte) (*pbmsg.Message, error) {
	id, err := message.GenID(nn.GetLocalNode().MessageIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.Bytes{
		Data: data,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_BYTES,
		RoutingType: pbmsg.RoutingType_DIRECT,
		MessageId:   id,
		Message:     buf,
		SrcId:       srcID,
		DestId:      key,
	}

	return msg, nil
}

// NewBroadcastBytesMessage creates a BYTES message that send arbitrary bytes to
// EVERY remote node in the network (not just neighbors)
func (nn *NNet) NewBroadcastBytesMessage(data, srcID []byte, routingType pbmsg.RoutingType) (*pbmsg.Message, error) {
	id, err := message.GenID(nn.GetLocalNode().MessageIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.Bytes{
		Data: data,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_BYTES,
		RoutingType: routingType,
		MessageId:   id,
		Message:     buf,
		SrcId:       srcID,
	}

	return msg, nil
}

// SendBytesDirectAsync sends bytes data to a remote node
func (nn *NNet) SendBytesDirectAsync(data []byte, remoteNode *node.RemoteNode) error {
	msg, err := nn.NewDirectBytesMessage(data)
	if err != nil {
		return err
	}

	return remoteNode.SendMessageAsync(msg)
}

// SendBytesDirectSync is the same as SendBytesDirectSyncWithTimeout but use
// default reply timeout in config.
func (nn *NNet) SendBytesDirectSync(data []byte, remoteNode *node.RemoteNode) ([]byte, *node.RemoteNode, error) {
	return nn.SendBytesDirectSyncWithTimeout(data, remoteNode, 0)
}

// SendBytesDirectSyncWithTimeout sends bytes data to a remote node, returns
// reply message, RemoteNode that sends the reply and error. Will also returns
// an error if doesn't receive a reply before timeout.
func (nn *NNet) SendBytesDirectSyncWithTimeout(data []byte, remoteNode *node.RemoteNode, replyTimeout time.Duration) ([]byte, *node.RemoteNode, error) {
	msg, err := nn.NewDirectBytesMessage(data)
	if err != nil {
		return nil, nil, err
	}

	reply, err := remoteNode.SendMessageSync(msg, replyTimeout)
	if err != nil {
		return nil, nil, err
	}

	replyBody := &pbmsg.Bytes{}
	err = proto.Unmarshal(reply.Msg.Message, replyBody)
	if err != nil {
		return nil, reply.RemoteNode, err
	}

	return replyBody.Data, reply.RemoteNode, nil
}

// SendBytesDirectReply is the same as SendBytesDirectAsync but with the
// replyToId field of the message set
func (nn *NNet) SendBytesDirectReply(replyToID, data []byte, remoteNode *node.RemoteNode) error {
	msg, err := nn.NewDirectBytesMessage(data)
	if err != nil {
		return err
	}

	msg.ReplyToId = replyToID

	return remoteNode.SendMessageAsync(msg)
}

// SendBytesRelayAsync sends bytes data to the remote node that has smallest
// distance to the key, returns if send success (which is true if successfully
// send message to at least one next hop), and aggregated error during message
// sending
func (nn *NNet) SendBytesRelayAsync(data, key []byte) (bool, error) {
	msg, err := nn.NewRelayBytesMessage(data, nn.GetLocalNode().Id, key)
	if err != nil {
		return false, err
	}

	return nn.SendMessageAsync(msg, pbmsg.RoutingType_RELAY)
}

// SendBytesRelaySync is the same as SendBytesRelaySync but use default reply
// timeout in config.
func (nn *NNet) SendBytesRelaySync(data, key []byte) ([]byte, []byte, error) {
	return nn.SendBytesRelaySyncWithTimeout(data, key, 0)
}

// SendBytesRelaySyncWithTimeout sends bytes data to the remote node that has
// smallest distance to the key, returns reply message, node ID who sends the
// reply, and aggregated error during message sending and receiving, will also
// returns error if doesn't receive any reply before timeout
func (nn *NNet) SendBytesRelaySyncWithTimeout(data, key []byte, replyTimeout time.Duration) ([]byte, []byte, error) {
	msg, err := nn.NewRelayBytesMessage(data, nn.GetLocalNode().Id, key)
	if err != nil {
		return nil, nil, err
	}

	reply, _, err := nn.SendMessageSync(msg, pbmsg.RoutingType_RELAY, replyTimeout)
	if err != nil {
		return nil, nil, err
	}

	replyBody := &pbmsg.Bytes{}
	err = proto.Unmarshal(reply.Message, replyBody)
	if err != nil {
		return nil, reply.SrcId, err
	}

	return replyBody.Data, reply.SrcId, nil
}

// SendBytesRelayReply is the same as SendBytesRelayAsync but with the
// replyToId field of the message set
func (nn *NNet) SendBytesRelayReply(replyToID, data, key []byte) (bool, error) {
	msg, err := nn.NewRelayBytesMessage(data, nn.GetLocalNode().Id, key)
	if err != nil {
		return false, err
	}

	msg.ReplyToId = replyToID

	return nn.SendMessageAsync(msg, pbmsg.RoutingType_RELAY)
}

// SendBytesBroadcastAsync sends bytes data to EVERY remote node in the network
// (not just neighbors), returns if send success (which is true if successfully
// send message to at least one next hop), and aggregated error during message
// sending
func (nn *NNet) SendBytesBroadcastAsync(data []byte, routingType pbmsg.RoutingType) (bool, error) {
	msg, err := nn.NewBroadcastBytesMessage(data, nn.GetLocalNode().Id, routingType)
	if err != nil {
		return false, err
	}

	return nn.SendMessageAsync(msg, routingType)
}

// SendBytesBroadcastSync is the same as SendBytesBroadcastSyncTimeout but use
// default reply timeout in config.
func (nn *NNet) SendBytesBroadcastSync(data []byte, routingType pbmsg.RoutingType) ([]byte, []byte, error) {
	return nn.SendBytesBroadcastSyncWithTimeout(data, routingType, 0)
}

// SendBytesBroadcastSyncWithTimeout sends bytes data to EVERY remote node in
// the network (not just neighbors), returns reply message, node ID who sends
// the reply, and aggregated error during message sending, will also returns
// error if doesn't receive any reply before timeout. Broadcast msg reply should
// be handled VERY carefully with some sort of sender identity verification,
// otherwise it may be used to DDoS attack sender with HUGE amplification
// factor.
func (nn *NNet) SendBytesBroadcastSyncWithTimeout(data []byte, routingType pbmsg.RoutingType, replyTimeout time.Duration) ([]byte, []byte, error) {
	msg, err := nn.NewBroadcastBytesMessage(data, nn.GetLocalNode().Id, routingType)
	if err != nil {
		return nil, nil, err
	}

	reply, _, err := nn.SendMessageSync(msg, routingType, replyTimeout)
	if err != nil {
		return nil, nil, err
	}

	replyBody := &pbmsg.Bytes{}
	err = proto.Unmarshal(reply.Message, replyBody)
	if err != nil {
		return nil, reply.SrcId, err
	}

	return replyBody.Data, reply.SrcId, nil
}

// SendBytesBroadcastReply is the same as SendBytesBroadcastAsync but with the
// replyToId field of the message set. Broadcast msg reply should be handled
// VERY carefully with some sort of sender identity verification, otherwise it
// may be used to DDoS attack sender with HUGE amplification factor. This should
// NOT be used unless msg src id is unknown.
func (nn *NNet) SendBytesBroadcastReply(replyToID, data []byte, routingType pbmsg.RoutingType) (bool, error) {
	msg, err := nn.NewBroadcastBytesMessage(data, nn.GetLocalNode().Id, routingType)
	if err != nil {
		return false, err
	}

	msg.ReplyToId = replyToID

	return nn.SendMessageAsync(msg, routingType)
}
