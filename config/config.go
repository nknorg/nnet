package config

import (
	"time"

	"github.com/imdario/mergo"
)

// Config is the configuration struct
type Config struct {
	Transport      string // which transport to use, e.g. tcp, udp, kcp
	Hostname       string // IP or domain name for remote node to connect to, e.g. 127.0.0.1, nkn.org. Empty string means remote nodes will fill it with your address they saw, which works if all nodes are not in the same local network or are all in the local network, but will cause problem if some nodes are in the same local network
	Port           uint16 // port to listen to incoming connections
	NodeIDBytes    uint32 // length of node id in bytes
	MessageIDBytes uint8  // MsgIDBytes is the length of message id in RandBytes

	Multiplexer        string // which multiplexer to use, e.g. smux, yamux
	NumStreamsToOpen   uint32 // number of streams to open per remote node
	NumStreamsToAccept uint32 // number of streams to accept per remote node

	LocalRxMsgChanLen              uint32        // Max number of msg that can be buffered per routing type
	LocalHandleMsgChanLen          uint32        // Max number of msg to be processed that can be buffered
	LocalRxMsgCacheExpiration      time.Duration // How long a received message id stays in cache before expiration
	LocalRxMsgCacheCleanupInterval time.Duration // How often to check and delete expired received message id

	RemoteRxMsgChanLen              uint32        // Max number of msg received that can be buffered
	RemoteTxMsgChanLen              uint32        // Max number of msg to be sent that can be buffered
	RemoteTxMsgCacheExpiration      time.Duration // How long a sent message id stays in cache before expiration
	RemoteTxMsgCacheCleanupInterval time.Duration // How often to check and delete expired sent message

	MaxMessageSize               uint32        // Max message size in bytes
	DefaultReplyTimeout          time.Duration // default timeout for receiving reply msg
	ReplyChanCleanupInterval     time.Duration // How often to check and delete expired reply chan
	MeasureRoundTripTimeInterval time.Duration // Time interval between measuring round trip time
	KeepAliveTimeout             time.Duration // Max idle time before considering node dead and closing connection
	DialTimeout                  time.Duration // Transport dial timeout

	OverlayLocalMsgChanLen uint32 // Max number of msg to be processed by local node that can be buffered

	MinNumSuccessors      uint32        // minimal number of successors of each chord node
	NumFingerSuccessors   uint32        // minimal number of successors of each finger table key
	NumSuccessorsFactor   uint32        // number of successors is max(this factor times the number of non empty finger table, MinNumSuccessors)
	BaseStabilizeInterval time.Duration // base stabilize interval
}

// DefaultConfig returns the default configurations
func DefaultConfig() *Config {
	defaultConfig := &Config{
		Transport:      "tcp",
		NodeIDBytes:    32,
		MessageIDBytes: 8,

		Multiplexer:        "smux",
		NumStreamsToOpen:   8,
		NumStreamsToAccept: 32,

		LocalRxMsgChanLen:              23333,
		LocalHandleMsgChanLen:          23333,
		LocalRxMsgCacheExpiration:      300 * time.Second,
		LocalRxMsgCacheCleanupInterval: 10 * time.Second,

		RemoteRxMsgChanLen:              2333,
		RemoteTxMsgChanLen:              2333,
		RemoteTxMsgCacheExpiration:      300 * time.Second,
		RemoteTxMsgCacheCleanupInterval: 10 * time.Second,

		MaxMessageSize:               20 * 1024 * 1024,
		DefaultReplyTimeout:          5 * time.Second,
		ReplyChanCleanupInterval:     1 * time.Second,
		MeasureRoundTripTimeInterval: 5 * time.Second,
		KeepAliveTimeout:             20 * time.Second,
		DialTimeout:                  5 * time.Second,

		OverlayLocalMsgChanLen: 23333,

		MinNumSuccessors:      8,
		NumFingerSuccessors:   3,
		NumSuccessorsFactor:   2,
		BaseStabilizeInterval: 2 * time.Second,
	}
	return defaultConfig
}

// MergedConfig returns a new Config that use fields in conf if provided,
// otherwise use default config
func MergedConfig(conf *Config) (*Config, error) {
	merged := DefaultConfig()
	err := mergo.Merge(merged, conf, mergo.WithOverride)
	if err != nil {
		return nil, err
	}
	return merged, nil
}
