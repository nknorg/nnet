package overlay

// NetworkWillStart is called right before network starts listening and handling
// messages. It can be used to do additional network setup like NAT traversal.
type NetworkWillStart struct {
	Func     func(Network) bool
	Priority int32
}

// NetworkStarted is called right after network starts listening and handling
// messages.
type NetworkStarted struct {
	Func     func(Network) bool
	Priority int32
}

// NetworkWillStop is called right before network stops listening and handling
// messages.
type NetworkWillStop struct {
	Func     func(Network) bool
	Priority int32
}

// NetworkStopped is called right after network stops listening and handling
// messages. It can be used to do some clean up, etc.
type NetworkStopped struct {
	Func     func(Network) bool
	Priority int32
}
