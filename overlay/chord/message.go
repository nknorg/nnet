package chord

import (
	"github.com/nknorg/nnet/message"
	"github.com/nknorg/nnet/node"
	pbmsg "github.com/nknorg/nnet/protobuf/message"
	pbnode "github.com/nknorg/nnet/protobuf/node"
	"google.golang.org/protobuf/proto"
)

// NewGetSuccAndPredMessage creates a GET_SUCC_AND_PRED message to get the
// successors and predecessor of a remote node
func NewGetSuccAndPredMessage(numSucc, numPred uint32, msgIDBytes uint8) (*pbmsg.Message, error) {
	id, err := message.GenID(msgIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.GetSuccAndPred{
		NumSucc: numSucc,
		NumPred: numPred,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_GET_SUCC_AND_PRED,
		RoutingType: pbmsg.RoutingType_DIRECT,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// NewGetSuccAndPredReply creates a GET_SUCC_AND_PRED reply to send successors
// and predecessor
func (c *Chord) NewGetSuccAndPredReply(replyToID []byte, successors, predecessors []*pbnode.Node) (*pbmsg.Message, error) {
	id, err := message.GenID(c.LocalNode.MessageIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.GetSuccAndPredReply{
		Successors:   successors,
		Predecessors: predecessors,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_GET_SUCC_AND_PRED,
		RoutingType: pbmsg.RoutingType_DIRECT,
		ReplyToId:   replyToID,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// NewFindSuccAndPredMessage creates a FIND_SUCC_AND_PRED message to find
// numSucc successors and numPred predecessors of a key
func (c *Chord) NewFindSuccAndPredMessage(key []byte, numSucc, numPred uint32) (*pbmsg.Message, error) {
	id, err := message.GenID(c.LocalNode.MessageIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.FindSuccAndPred{
		Key:     key,
		NumSucc: numSucc,
		NumPred: numPred,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_FIND_SUCC_AND_PRED,
		RoutingType: pbmsg.RoutingType_DIRECT,
		MessageId:   id,
		Message:     buf,
		DestId:      key,
	}

	return msg, nil
}

// NewFindSuccAndPredReply creates a FIND_SUCC_AND_PRED reply to send successors
// and predecessors
func (c *Chord) NewFindSuccAndPredReply(replyToID []byte, successors, predecessors []*pbnode.Node) (*pbmsg.Message, error) {
	id, err := message.GenID(c.LocalNode.MessageIDBytes)
	if err != nil {
		return nil, err
	}

	msgBody := &pbmsg.FindSuccAndPredReply{
		Successors:   successors,
		Predecessors: predecessors,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pbmsg.Message{
		MessageType: pbmsg.MessageType_FIND_SUCC_AND_PRED,
		RoutingType: pbmsg.RoutingType_DIRECT,
		ReplyToId:   replyToID,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// handleRemoteMessage handles a remote message and returns if it should be
// passed through to local node and error
func (c *Chord) handleRemoteMessage(remoteMsg *node.RemoteMessage) (bool, error) {
	if remoteMsg.RemoteNode == nil {
		return true, nil
	}

	switch remoteMsg.Msg.MessageType {
	case pbmsg.MessageType_GET_SUCC_AND_PRED:
		msgBody := &pbmsg.GetSuccAndPred{}
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

		replyMsg, err := c.NewGetSuccAndPredReply(remoteMsg.Msg.MessageId, succs, preds)
		if err != nil {
			return false, err
		}

		err = remoteMsg.RemoteNode.SendMessageAsync(replyMsg)
		if err != nil {
			return false, err
		}

	case pbmsg.MessageType_FIND_SUCC_AND_PRED:
		msgBody := &pbmsg.FindSuccAndPred{}
		err := proto.Unmarshal(remoteMsg.Msg.Message, msgBody)
		if err != nil {
			return false, err
		}

		// TODO: prevent unbounded number of goroutines
		go func() {
			succs, preds, err := c.FindSuccAndPred(msgBody.Key, msgBody.NumSucc, msgBody.NumPred)
			if err != nil {
				return
			}

			replyMsg, err := c.NewFindSuccAndPredReply(remoteMsg.Msg.MessageId, succs, preds)
			if err != nil {
				return
			}

			err = remoteMsg.RemoteNode.SendMessageAsync(replyMsg)
			if err != nil {
				return
			}
		}()

	default:
		return true, nil
	}

	return false, nil
}
