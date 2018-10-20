package nnet

import (
	"github.com/gogo/protobuf/proto"
	"github.com/nknorg/nnet/message"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/protobuf"
)

// NewDirectBytesMessage creates a BYTES message that send arbitrary bytes to a
// given remote node
func NewDirectBytesMessage(data []byte) (*protobuf.Message, error) {
	id, err := message.GenID()
	if err != nil {
		return nil, err
	}

	msgBody := &protobuf.Bytes{
		Data: data,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		MessageType: protobuf.BYTES,
		RoutingType: protobuf.DIRECT,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// NewRelayBytesMessage creates a BYTES message that send arbitrary bytes to the
// remote node that has the smallest distance to a given key
func NewRelayBytesMessage(data, srcID, key []byte) (*protobuf.Message, error) {
	id, err := message.GenID()
	if err != nil {
		return nil, err
	}

	msgBody := &protobuf.Bytes{
		Data: data,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		MessageType: protobuf.BYTES,
		RoutingType: protobuf.RELAY,
		MessageId:   id,
		Message:     buf,
		SrcId:       srcID,
		DestId:      key,
	}

	return msg, nil
}

// NewBroadcastBytesMessage creates a BYTES message that send arbitrary bytes to
// EVERY remote node in the network (not just neighbors)
func NewBroadcastBytesMessage(data, srcID []byte) (*protobuf.Message, error) {
	id, err := message.GenID()
	if err != nil {
		return nil, err
	}

	msgBody := &protobuf.Bytes{
		Data: data,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		MessageType: protobuf.BYTES,
		RoutingType: protobuf.BROADCAST,
		MessageId:   id,
		Message:     buf,
		SrcId:       srcID,
	}

	return msg, nil
}

// SendBytesDirectAsync sends bytes data to a remote node
func (nn *NNet) SendBytesDirectAsync(data []byte, remoteNode *node.RemoteNode) error {
	msg, err := NewDirectBytesMessage(data)
	if err != nil {
		return err
	}

	return remoteNode.SendMessageAsync(msg)
}

// SendBytesDirectSync sends bytes data to a remote node, returns reply message
// or error if timeout
func (nn *NNet) SendBytesDirectSync(data []byte, remoteNode *node.RemoteNode) (*protobuf.Message, error) {
	msg, err := NewDirectBytesMessage(data)
	if err != nil {
		return nil, err
	}

	reply, err := remoteNode.SendMessageSync(msg)
	if err != nil {
		return nil, err
	}

	return reply.Msg, nil
}

// SendBytesRelayAsync sends bytes data to the remote node that has smallest
// distance to the key, returns if send success (which is true if successfully
// send message to at least one next hop), and aggregated error during message
// sending
func (nn *NNet) SendBytesRelayAsync(data, key []byte) (bool, error) {
	msg, err := NewRelayBytesMessage(data, nn.LocalNode.Id, key)
	if err != nil {
		return false, err
	}

	return nn.Overlay.SendMessageAsync(msg, protobuf.RELAY)
}

// SendBytesRelaySync sends bytes data to the remote node that has smallest
// distance to the key, returns reply message, if send success (which is true if
// successfully send message to at least one next hop), and aggregated error
// during message sending, will also returns error if doesn't receive any reply
// before timeout
func (nn *NNet) SendBytesRelaySync(data, key []byte) (*protobuf.Message, bool, error) {
	msg, err := NewRelayBytesMessage(data, nn.LocalNode.Id, key)
	if err != nil {
		return nil, false, err
	}

	return nn.Overlay.SendMessageSync(msg, protobuf.RELAY)
}

// SendBytesBroadcastAsync sends bytes data to EVERY remote node in the network
// (not just neighbors), returns if send success (which is true if successfully
// send message to at least one next hop), and aggregated error during message
// sending
func (nn *NNet) SendBytesBroadcastAsync(data []byte) (bool, error) {
	msg, err := NewBroadcastBytesMessage(data, nn.LocalNode.Id)
	if err != nil {
		return false, err
	}

	return nn.Overlay.SendMessageAsync(msg, protobuf.BROADCAST)
}

// SendBytesBroadcastSync sends bytes data to EVERY remote node in the network
// (not just neighbors), returns reply message, if send success (which is true
// if successfully send message to at least one next hop), and aggregated error
// during message sending, will also returns error if doesn't receive any reply
// before timeout
func (nn *NNet) SendBytesBroadcastSync(data []byte) (*protobuf.Message, bool, error) {
	msg, err := NewBroadcastBytesMessage(data, nn.LocalNode.Id)
	if err != nil {
		return nil, false, err
	}

	return nn.Overlay.SendMessageSync(msg, protobuf.BROADCAST)
}
