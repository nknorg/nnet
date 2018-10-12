package chord

import (
	"github.com/gogo/protobuf/proto"
	"github.com/nknorg/nnet/message"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/protobuf"
)

// NewGetSuccAndPredMessage creates a GET_SUCC_AND_PRED message to get the
// successors and predecessor of a remote node
func NewGetSuccAndPredMessage(numSucc, numPred uint32) (*protobuf.Message, error) {
	id, err := message.GenID()
	if err != nil {
		return nil, err
	}

	msgBody := &protobuf.GetSuccAndPred{
		NumSucc: numSucc,
		NumPred: numPred,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		MessageType: protobuf.GET_SUCC_AND_PRED,
		RoutingType: protobuf.DIRECT,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// NewGetSuccAndPredReply creates a GET_SUCC_AND_PRED reply to send successors
// and predecessor
func NewGetSuccAndPredReply(replyToID []byte, successors, predecessors []*protobuf.Node) (*protobuf.Message, error) {
	msgBody := &protobuf.GetSuccAndPredReply{
		Successors:   successors,
		Predecessors: predecessors,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		MessageType: protobuf.GET_SUCC_AND_PRED,
		RoutingType: protobuf.DIRECT,
		ReplyToId:   replyToID,
		Message:     buf,
	}

	return msg, nil
}

// NewFindSuccessorsMessage creates a FIND_SUCCESSORS message to find numNodes
// successors of a key
func NewFindSuccessorsMessage(key []byte, numSucc uint32) (*protobuf.Message, error) {
	id, err := message.GenID()
	if err != nil {
		return nil, err
	}

	msgBody := &protobuf.FindSuccessors{
		Key:     key,
		NumSucc: numSucc,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		MessageType: protobuf.FIND_SUCCESSORS,
		RoutingType: protobuf.DIRECT,
		MessageId:   id,
		Message:     buf,
		DestId:      key,
	}

	return msg, nil
}

// NewFindSuccessorsReply creates a FIND_SUCCESSORS reply to send successors
func NewFindSuccessorsReply(replyToID []byte, successors []*protobuf.Node) (*protobuf.Message, error) {
	msgBody := &protobuf.FindSuccessorsReply{
		Successors: successors,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		MessageType: protobuf.FIND_SUCCESSORS,
		RoutingType: protobuf.DIRECT,
		ReplyToId:   replyToID,
		Message:     buf,
	}

	return msg, nil
}

// handleRemoteMessage handles a remote message and returns if it should be
// passed through to local node and error
func (c *Chord) handleRemoteMessage(remoteMsg *node.RemoteMessage) (bool, error) {
	switch remoteMsg.Msg.MessageType {
	case protobuf.GET_SUCC_AND_PRED:
		msgBody := &protobuf.GetSuccAndPred{}
		err := proto.Unmarshal(remoteMsg.Msg.Message, msgBody)
		if err != nil {
			return false, err
		}

		succs := c.successors.ToProtoNodeList(true)
		if succs != nil && uint32(len(succs)) > msgBody.NumSucc {
			succs = succs[:msgBody.NumSucc]
		}

		preds := c.predecessors.ToProtoNodeList(true)
		if preds != nil && uint32(len(preds)) > msgBody.NumPred {
			preds = preds[:msgBody.NumPred]
		}

		replyMsg, err := NewGetSuccAndPredReply(remoteMsg.Msg.MessageId, succs, preds)
		if err != nil {
			return false, err
		}

		err = remoteMsg.RemoteNode.SendMessageAsync(replyMsg)
		if err != nil {
			return false, err
		}

	case protobuf.FIND_SUCCESSORS:
		msgBody := &protobuf.FindSuccessors{}
		err := proto.Unmarshal(remoteMsg.Msg.Message, msgBody)
		if err != nil {
			return false, err
		}

		succs, err := c.FindSuccessors(msgBody.Key, msgBody.NumSucc)
		if err != nil {
			return false, err
		}

		replyMsg, err := NewFindSuccessorsReply(remoteMsg.Msg.MessageId, succs)
		if err != nil {
			return false, err
		}

		err = remoteMsg.RemoteNode.SendMessageAsync(replyMsg)
		if err != nil {
			return false, err
		}

	case protobuf.FIND_PREDECESSOR:
	default:
		return true, nil
	}

	return false, nil
}
