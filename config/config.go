package config

import (
	"time"

	"github.com/imdario/mergo"
)

// Config is the exposed configurations
type Config struct {
	Transport            string        // which transport to use, e.g. tcp, udp, kcp
	Port                 uint16        // port to listen to incoming connections
	NodeIDBytes          uint32        // length of node id in bytes
	MinNumSuccessors     uint32        // minimal number of successors of each chord node
	MinNumPredecessors   uint32        // minimal number of predecessors of each chord node
	NumFingerSuccessors  uint32        // minimal number of successors of each finger table key
	MinStabilizeInterval time.Duration // minimal stabilize interval
	MaxStabilizeInterval time.Duration // maximal stabilize interval
	DialTimeout          time.Duration // dial timeout for tcp and udp
}

// DefaultConfig returns the default configurations
func DefaultConfig() *Config {
	defaultConfig := &Config{
		Transport:            "tcp",
		NodeIDBytes:          32,
		MinNumSuccessors:     8,
		MinNumPredecessors:   1,
		NumFingerSuccessors:  1,
		MinStabilizeInterval: 1 * time.Second,
		MaxStabilizeInterval: 2 * time.Second,
		DialTimeout:          5 * time.Second,
	}
	return defaultConfig
}

// MergedConfig returns a new Config that use fields in conf if provided,
// otherwise use default config
func MergedConfig(conf Config) (*Config, error) {
	merged := DefaultConfig()
	err := mergo.Merge(merged, conf, mergo.WithOverride)
	if err != nil {
		return nil, err
	}
	return merged, nil
}
