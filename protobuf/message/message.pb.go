// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: protobuf/message/message.proto

package message

import (
	node "github.com/nknorg/nnet/protobuf/node"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RoutingType int32

const (
	RoutingType_DIRECT         RoutingType = 0
	RoutingType_RELAY          RoutingType = 1
	RoutingType_BROADCAST_PUSH RoutingType = 2
	RoutingType_BROADCAST_PULL RoutingType = 3 // to be implemented
	RoutingType_BROADCAST_TREE RoutingType = 4
)

// Enum value maps for RoutingType.
var (
	RoutingType_name = map[int32]string{
		0: "DIRECT",
		1: "RELAY",
		2: "BROADCAST_PUSH",
		3: "BROADCAST_PULL",
		4: "BROADCAST_TREE",
	}
	RoutingType_value = map[string]int32{
		"DIRECT":         0,
		"RELAY":          1,
		"BROADCAST_PUSH": 2,
		"BROADCAST_PULL": 3,
		"BROADCAST_TREE": 4,
	}
)

func (x RoutingType) Enum() *RoutingType {
	p := new(RoutingType)
	*p = x
	return p
}

func (x RoutingType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RoutingType) Descriptor() protoreflect.EnumDescriptor {
	return file_protobuf_message_message_proto_enumTypes[0].Descriptor()
}

func (RoutingType) Type() protoreflect.EnumType {
	return &file_protobuf_message_message_proto_enumTypes[0]
}

func (x RoutingType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RoutingType.Descriptor instead.
func (RoutingType) EnumDescriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{0}
}

type MessageType int32

const (
	// Node message
	MessageType_PING          MessageType = 0
	MessageType_EXCHANGE_NODE MessageType = 1
	MessageType_STOP          MessageType = 2
	// Chord message
	MessageType_GET_SUCC_AND_PRED  MessageType = 3
	MessageType_FIND_SUCC_AND_PRED MessageType = 4
	// Message that contains any bytes
	MessageType_BYTES MessageType = 5
)

// Enum value maps for MessageType.
var (
	MessageType_name = map[int32]string{
		0: "PING",
		1: "EXCHANGE_NODE",
		2: "STOP",
		3: "GET_SUCC_AND_PRED",
		4: "FIND_SUCC_AND_PRED",
		5: "BYTES",
	}
	MessageType_value = map[string]int32{
		"PING":               0,
		"EXCHANGE_NODE":      1,
		"STOP":               2,
		"GET_SUCC_AND_PRED":  3,
		"FIND_SUCC_AND_PRED": 4,
		"BYTES":              5,
	}
)

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}

func (x MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_protobuf_message_message_proto_enumTypes[1].Descriptor()
}

func (MessageType) Type() protoreflect.EnumType {
	return &file_protobuf_message_message_proto_enumTypes[1]
}

func (x MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageType.Descriptor instead.
func (MessageType) EnumDescriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{1}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoutingType RoutingType `protobuf:"varint,1,opt,name=routing_type,json=routingType,proto3,enum=RoutingType" json:"routing_type,omitempty"`
	MessageType MessageType `protobuf:"varint,2,opt,name=message_type,json=messageType,proto3,enum=MessageType" json:"message_type,omitempty"`
	Message     []byte      `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	MessageId   []byte      `protobuf:"bytes,4,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	ReplyToId   []byte      `protobuf:"bytes,5,opt,name=reply_to_id,json=replyToId,proto3" json:"reply_to_id,omitempty"`
	SrcId       []byte      `protobuf:"bytes,6,opt,name=src_id,json=srcId,proto3" json:"src_id,omitempty"`
	DestId      []byte      `protobuf:"bytes,7,opt,name=dest_id,json=destId,proto3" json:"dest_id,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetRoutingType() RoutingType {
	if x != nil {
		return x.RoutingType
	}
	return RoutingType_DIRECT
}

func (x *Message) GetMessageType() MessageType {
	if x != nil {
		return x.MessageType
	}
	return MessageType_PING
}

func (x *Message) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *Message) GetMessageId() []byte {
	if x != nil {
		return x.MessageId
	}
	return nil
}

func (x *Message) GetReplyToId() []byte {
	if x != nil {
		return x.ReplyToId
	}
	return nil
}

func (x *Message) GetSrcId() []byte {
	if x != nil {
		return x.SrcId
	}
	return nil
}

func (x *Message) GetDestId() []byte {
	if x != nil {
		return x.DestId
	}
	return nil
}

type Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Ping) Reset() {
	*x = Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ping) ProtoMessage() {}

func (x *Ping) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ping.ProtoReflect.Descriptor instead.
func (*Ping) Descriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{1}
}

type PingReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PingReply) Reset() {
	*x = PingReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingReply) ProtoMessage() {}

func (x *PingReply) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingReply.ProtoReflect.Descriptor instead.
func (*PingReply) Descriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{2}
}

type ExchangeNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node *node.Node `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
}

func (x *ExchangeNode) Reset() {
	*x = ExchangeNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExchangeNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExchangeNode) ProtoMessage() {}

func (x *ExchangeNode) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExchangeNode.ProtoReflect.Descriptor instead.
func (*ExchangeNode) Descriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{3}
}

func (x *ExchangeNode) GetNode() *node.Node {
	if x != nil {
		return x.Node
	}
	return nil
}

type ExchangeNodeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node *node.Node `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
}

func (x *ExchangeNodeReply) Reset() {
	*x = ExchangeNodeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_message_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExchangeNodeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExchangeNodeReply) ProtoMessage() {}

func (x *ExchangeNodeReply) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_message_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExchangeNodeReply.ProtoReflect.Descriptor instead.
func (*ExchangeNodeReply) Descriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{4}
}

func (x *ExchangeNodeReply) GetNode() *node.Node {
	if x != nil {
		return x.Node
	}
	return nil
}

type Stop struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Stop) Reset() {
	*x = Stop{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_message_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stop) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stop) ProtoMessage() {}

func (x *Stop) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_message_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stop.ProtoReflect.Descriptor instead.
func (*Stop) Descriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{5}
}

type GetSuccAndPred struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NumSucc uint32 `protobuf:"varint,1,opt,name=num_succ,json=numSucc,proto3" json:"num_succ,omitempty"`
	NumPred uint32 `protobuf:"varint,2,opt,name=num_pred,json=numPred,proto3" json:"num_pred,omitempty"`
}

func (x *GetSuccAndPred) Reset() {
	*x = GetSuccAndPred{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_message_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSuccAndPred) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSuccAndPred) ProtoMessage() {}

func (x *GetSuccAndPred) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_message_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSuccAndPred.ProtoReflect.Descriptor instead.
func (*GetSuccAndPred) Descriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{6}
}

func (x *GetSuccAndPred) GetNumSucc() uint32 {
	if x != nil {
		return x.NumSucc
	}
	return 0
}

func (x *GetSuccAndPred) GetNumPred() uint32 {
	if x != nil {
		return x.NumPred
	}
	return 0
}

type GetSuccAndPredReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Successors   []*node.Node `protobuf:"bytes,1,rep,name=successors,proto3" json:"successors,omitempty"`
	Predecessors []*node.Node `protobuf:"bytes,2,rep,name=predecessors,proto3" json:"predecessors,omitempty"`
}

func (x *GetSuccAndPredReply) Reset() {
	*x = GetSuccAndPredReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_message_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSuccAndPredReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSuccAndPredReply) ProtoMessage() {}

func (x *GetSuccAndPredReply) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_message_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSuccAndPredReply.ProtoReflect.Descriptor instead.
func (*GetSuccAndPredReply) Descriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{7}
}

func (x *GetSuccAndPredReply) GetSuccessors() []*node.Node {
	if x != nil {
		return x.Successors
	}
	return nil
}

func (x *GetSuccAndPredReply) GetPredecessors() []*node.Node {
	if x != nil {
		return x.Predecessors
	}
	return nil
}

type FindSuccAndPred struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key     []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	NumSucc uint32 `protobuf:"varint,2,opt,name=num_succ,json=numSucc,proto3" json:"num_succ,omitempty"`
	NumPred uint32 `protobuf:"varint,3,opt,name=num_pred,json=numPred,proto3" json:"num_pred,omitempty"`
}

func (x *FindSuccAndPred) Reset() {
	*x = FindSuccAndPred{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_message_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindSuccAndPred) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindSuccAndPred) ProtoMessage() {}

func (x *FindSuccAndPred) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_message_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindSuccAndPred.ProtoReflect.Descriptor instead.
func (*FindSuccAndPred) Descriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{8}
}

func (x *FindSuccAndPred) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *FindSuccAndPred) GetNumSucc() uint32 {
	if x != nil {
		return x.NumSucc
	}
	return 0
}

func (x *FindSuccAndPred) GetNumPred() uint32 {
	if x != nil {
		return x.NumPred
	}
	return 0
}

type FindSuccAndPredReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Successors   []*node.Node `protobuf:"bytes,1,rep,name=successors,proto3" json:"successors,omitempty"`
	Predecessors []*node.Node `protobuf:"bytes,2,rep,name=predecessors,proto3" json:"predecessors,omitempty"`
}

func (x *FindSuccAndPredReply) Reset() {
	*x = FindSuccAndPredReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_message_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindSuccAndPredReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindSuccAndPredReply) ProtoMessage() {}

func (x *FindSuccAndPredReply) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_message_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindSuccAndPredReply.ProtoReflect.Descriptor instead.
func (*FindSuccAndPredReply) Descriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{9}
}

func (x *FindSuccAndPredReply) GetSuccessors() []*node.Node {
	if x != nil {
		return x.Successors
	}
	return nil
}

func (x *FindSuccAndPredReply) GetPredecessors() []*node.Node {
	if x != nil {
		return x.Predecessors
	}
	return nil
}

type Bytes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Bytes) Reset() {
	*x = Bytes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_message_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bytes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bytes) ProtoMessage() {}

func (x *Bytes) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_message_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bytes.ProtoReflect.Descriptor instead.
func (*Bytes) Descriptor() ([]byte, []int) {
	return file_protobuf_message_message_proto_rawDescGZIP(), []int{10}
}

func (x *Bytes) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_protobuf_message_message_proto protoreflect.FileDescriptor

var file_protobuf_message_message_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2f,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf4, 0x01, 0x0a, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2f, 0x0a, 0x0c, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e,
	0x67, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x52,
	0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x72, 0x6f, 0x75, 0x74,
	0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2f, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49,
	0x64, 0x12, 0x1e, 0x0a, 0x0b, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x5f, 0x74, 0x6f, 0x5f, 0x69, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x49,
	0x64, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x72, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x73, 0x72, 0x63, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x64, 0x65, 0x73, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x64, 0x65, 0x73, 0x74, 0x49,
	0x64, 0x22, 0x06, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x22, 0x0b, 0x0a, 0x09, 0x50, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x29, 0x0a, 0x0c, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x19, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6e, 0x6f, 0x64,
	0x65, 0x22, 0x2e, 0x0a, 0x11, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x4e, 0x6f, 0x64,
	0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x19, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6e, 0x6f, 0x64,
	0x65, 0x22, 0x06, 0x0a, 0x04, 0x53, 0x74, 0x6f, 0x70, 0x22, 0x46, 0x0a, 0x0e, 0x47, 0x65, 0x74,
	0x53, 0x75, 0x63, 0x63, 0x41, 0x6e, 0x64, 0x50, 0x72, 0x65, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6e,
	0x75, 0x6d, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6e,
	0x75, 0x6d, 0x53, 0x75, 0x63, 0x63, 0x12, 0x19, 0x0a, 0x08, 0x6e, 0x75, 0x6d, 0x5f, 0x70, 0x72,
	0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6e, 0x75, 0x6d, 0x50, 0x72, 0x65,
	0x64, 0x22, 0x67, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x53, 0x75, 0x63, 0x63, 0x41, 0x6e, 0x64, 0x50,
	0x72, 0x65, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x25, 0x0a, 0x0a, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x6f, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x4e,
	0x6f, 0x64, 0x65, 0x52, 0x0a, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73, 0x12,
	0x29, 0x0a, 0x0c, 0x70, 0x72, 0x65, 0x64, 0x65, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x0c, 0x70, 0x72,
	0x65, 0x64, 0x65, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73, 0x22, 0x59, 0x0a, 0x0f, 0x46, 0x69,
	0x6e, 0x64, 0x53, 0x75, 0x63, 0x63, 0x41, 0x6e, 0x64, 0x50, 0x72, 0x65, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x19, 0x0a, 0x08, 0x6e, 0x75, 0x6d, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x07, 0x6e, 0x75, 0x6d, 0x53, 0x75, 0x63, 0x63, 0x12, 0x19, 0x0a, 0x08, 0x6e, 0x75,
	0x6d, 0x5f, 0x70, 0x72, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6e, 0x75,
	0x6d, 0x50, 0x72, 0x65, 0x64, 0x22, 0x68, 0x0a, 0x14, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x75, 0x63,
	0x63, 0x41, 0x6e, 0x64, 0x50, 0x72, 0x65, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x25, 0x0a,
	0x0a, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x05, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x0a, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x6f, 0x72, 0x73, 0x12, 0x29, 0x0a, 0x0c, 0x70, 0x72, 0x65, 0x64, 0x65, 0x63, 0x65, 0x73,
	0x73, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x4e, 0x6f, 0x64,
	0x65, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x64, 0x65, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73, 0x22,
	0x1b, 0x0a, 0x05, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x60, 0x0a, 0x0b,
	0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x44,
	0x49, 0x52, 0x45, 0x43, 0x54, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x52, 0x45, 0x4c, 0x41, 0x59,
	0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x42, 0x52, 0x4f, 0x41, 0x44, 0x43, 0x41, 0x53, 0x54, 0x5f,
	0x50, 0x55, 0x53, 0x48, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x42, 0x52, 0x4f, 0x41, 0x44, 0x43,
	0x41, 0x53, 0x54, 0x5f, 0x50, 0x55, 0x4c, 0x4c, 0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e, 0x42, 0x52,
	0x4f, 0x41, 0x44, 0x43, 0x41, 0x53, 0x54, 0x5f, 0x54, 0x52, 0x45, 0x45, 0x10, 0x04, 0x2a, 0x6e,
	0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a,
	0x04, 0x50, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x45, 0x58, 0x43, 0x48, 0x41,
	0x4e, 0x47, 0x45, 0x5f, 0x4e, 0x4f, 0x44, 0x45, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x54,
	0x4f, 0x50, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x45, 0x54, 0x5f, 0x53, 0x55, 0x43, 0x43,
	0x5f, 0x41, 0x4e, 0x44, 0x5f, 0x50, 0x52, 0x45, 0x44, 0x10, 0x03, 0x12, 0x16, 0x0a, 0x12, 0x46,
	0x49, 0x4e, 0x44, 0x5f, 0x53, 0x55, 0x43, 0x43, 0x5f, 0x41, 0x4e, 0x44, 0x5f, 0x50, 0x52, 0x45,
	0x44, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x42, 0x59, 0x54, 0x45, 0x53, 0x10, 0x05, 0x42, 0x29,
	0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x6b, 0x6e,
	0x6f, 0x72, 0x67, 0x2f, 0x6e, 0x6e, 0x65, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_protobuf_message_message_proto_rawDescOnce sync.Once
	file_protobuf_message_message_proto_rawDescData = file_protobuf_message_message_proto_rawDesc
)

func file_protobuf_message_message_proto_rawDescGZIP() []byte {
	file_protobuf_message_message_proto_rawDescOnce.Do(func() {
		file_protobuf_message_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_message_message_proto_rawDescData)
	})
	return file_protobuf_message_message_proto_rawDescData
}

var file_protobuf_message_message_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_protobuf_message_message_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_protobuf_message_message_proto_goTypes = []any{
	(RoutingType)(0),             // 0: RoutingType
	(MessageType)(0),             // 1: MessageType
	(*Message)(nil),              // 2: Message
	(*Ping)(nil),                 // 3: Ping
	(*PingReply)(nil),            // 4: PingReply
	(*ExchangeNode)(nil),         // 5: ExchangeNode
	(*ExchangeNodeReply)(nil),    // 6: ExchangeNodeReply
	(*Stop)(nil),                 // 7: Stop
	(*GetSuccAndPred)(nil),       // 8: GetSuccAndPred
	(*GetSuccAndPredReply)(nil),  // 9: GetSuccAndPredReply
	(*FindSuccAndPred)(nil),      // 10: FindSuccAndPred
	(*FindSuccAndPredReply)(nil), // 11: FindSuccAndPredReply
	(*Bytes)(nil),                // 12: Bytes
	(*node.Node)(nil),            // 13: Node
}
var file_protobuf_message_message_proto_depIdxs = []int32{
	0,  // 0: Message.routing_type:type_name -> RoutingType
	1,  // 1: Message.message_type:type_name -> MessageType
	13, // 2: ExchangeNode.node:type_name -> Node
	13, // 3: ExchangeNodeReply.node:type_name -> Node
	13, // 4: GetSuccAndPredReply.successors:type_name -> Node
	13, // 5: GetSuccAndPredReply.predecessors:type_name -> Node
	13, // 6: FindSuccAndPredReply.successors:type_name -> Node
	13, // 7: FindSuccAndPredReply.predecessors:type_name -> Node
	8,  // [8:8] is the sub-list for method output_type
	8,  // [8:8] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_protobuf_message_message_proto_init() }
func file_protobuf_message_message_proto_init() {
	if File_protobuf_message_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobuf_message_message_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_message_message_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Ping); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_message_message_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*PingReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_message_message_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ExchangeNode); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_message_message_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ExchangeNodeReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_message_message_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*Stop); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_message_message_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*GetSuccAndPred); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_message_message_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*GetSuccAndPredReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_message_message_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*FindSuccAndPred); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_message_message_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*FindSuccAndPredReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_message_message_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*Bytes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protobuf_message_message_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protobuf_message_message_proto_goTypes,
		DependencyIndexes: file_protobuf_message_message_proto_depIdxs,
		EnumInfos:         file_protobuf_message_message_proto_enumTypes,
		MessageInfos:      file_protobuf_message_message_proto_msgTypes,
	}.Build()
	File_protobuf_message_message_proto = out.File
	file_protobuf_message_message_proto_rawDesc = nil
	file_protobuf_message_message_proto_goTypes = nil
	file_protobuf_message_message_proto_depIdxs = nil
}