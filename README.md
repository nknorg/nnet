# nnet: a fast, scalable, and easy-to-use p2p overlay network stack

## Introduction

**nnet** is a fast, scalable, and developer-friendly network stack/layer for
decentralized/distributed systems. It handles message delivery, topology
maintenance, etc in a highly efficient and scalable way, while providing enough
control and flexibility through developer-friendly messaging interface and
powerful middleware architecture.

## Features

* nnet uses a **modular and layered overlay architecture**. By default an improved and much more reliable version of Chord DHT protocol is used to maintain a scalable overlay topology, while other topologies like Kademlia can be easily added by implementing a few overlay interfaces.
* Highly efficient messaging implementation that is able to send, receive and handle ~**75k messages/s** or ~**1 GB/s** of messages on a 2-core personal laptop.
* Extremely easy to use message sending interface with **both async and sync message flow** (block until reply). Message reply in sync mode are handled efficiently and automatically, you just need to provide the content.
* Deliver message to any node in the network (**not just the nodes you are directly connected to**) reliably and efficiently in at most log_2(N) hops (w.h.p) where N is the total number of nodes in the network.
* **Novel and highly efficient** message broadcasting algorithm with exact once (or K-times where K is something adjustable) message sending that achieves **optimal throughput and near-optimal latency**. This is done by sending message through the spanning tree constructed by utilizing the Chord topology.
* **Powerful and flexible middleware architecture** that allows you to easily interact with node/network lifecycle and events like topology change, routing, message sending and delivery, etc. **Applying a middleware is as simple as providing a function**.
* Flexible **transport-aware address scheme**. Each node can choose its own transport protocol to listen to, and nodes with different transport protocol can communicate with each other transparently. nnet supports TCP and [KCP](https://github.com/skywind3000/kcp/blob/master/README.en.md) transport layer by default. Other transport layers can be easily supported by implementing a few interfaces.
* **Modular and extensible router architecture**. Implementing a new routing algorithm is as simple as adding a router that implements a few router interfaces.
* Only **a fixed number of goroutines and connections** will be created given network size, and the number can be changed easily by changing the number of concurrent workers.
* **NAT traversal** (UPnP and NAT-PMP) using middleware.
* Use protocol buffers for message serialization/deserialization to support cross platform and backward/forward compatibility.
* Provide your own logger for logging integration with your application.

Coming soon:

- [x] Latency measurement
- [x] Proximity routing
- [ ] Push-pull broadcasting for large messages
- [ ] Efficient Pub/Sub
- [ ] Test cases

## Usage

### Requirements:

* Go 1.10+

### Install

```bash
go get -u -d github.com/nknorg/nnet
```

It's recommended to use package management tools such as
[glide](https://github.com/Masterminds/glide) to ensure reproducibility by
freezing dependencies at exact version/commit. If you have glide installed, you
can run `glide install` in `$GOPATH/src/github.com/nknorg/nnet/` directory to
use exactly the same version of dependencies as we do.

### Basic

Assuming you have imported `github.com/nknorg/nnet`, The bare minimal way to
create a nnet node is simply

```go
nn, err := nnet.NewNNet(nil, nil)
```

This will create a nnet node with random ID and default configuration (listen to
a random port, etc). Starting the node is as simple as

```go
err = nn.Start(true) // or false if joining rather than creating
```

To join an existing network, simply call `Join` after node has been started

```go
err = nn.Join("<protocol>://<ip>:<port>")
```

Put them together, a local nnet cluster with 10 nodes can be created by just a
few lines of code.

```go
nnets := make([]*nnet.NNet, 10)
for i := 0; i < len(nnets); i++ {
  nnets[i], _ = nnet.NewNNet(nil, nil) // error is omitted here for simplicity
  nnets[i].Start(i == 0) // error is omitted here for simplicity
  if i > 0 {
    nnets[i].Join(nnets[0].GetLocalNode().Addr) // error is omitted here for simplicity
  }
}
```

A complete basic example can be found at
[examples/basic/main.go](examples/basic/main.go). You can run it by

```bash
go run $GOPATH/src/github.com/nknorg/nnet/examples/basic/main.go
```

### Middleware

nnet uses middleware as a powerful and flexible architecture to interact with
node/network lifecycle and events such as topology change, routing, message
sending and delivery, etc. A middleware is just a function that, once hooked,
will be called when certain event occurs. Applying a middleware is as simple as
calling `ApplyMiddleware` function (or `MustApplyMiddleware` which will panic if
an error occurs). Let's apply a `node.RemoteNodeReady` middleware that will be
called when a remote node is connected with the local node and exchanged node
info

```go
err = nn.ApplyMiddleware(node.RemoteNodeReady(func(remoteNode *node.RemoteNode) bool {
  fmt.Printf("Remote node ready: %v", remoteNode)
  return true
}))
```

Different middleware types take different arguments and have different return
types, but they all share one thing in common: one of their returned value is a
boolean indicating whether we should call the next middleware (of the same
type). Multiple middleware of the same type will be called in the order added.
They form a pipeline such that each one can respond to the event with some data,
modify data that will be passed to the rest middleware, or decides to stop the
data flow. For example, we can randomly reject remote nodes

```go
nn.MustApplyMiddleware(node.RemoteNodeConnected(func(remoteNode *node.RemoteNode) bool {
  if rand.Float64() < 0.23333 {
    remoteNode.Stop(errors.New("YOU ARE UNLUCKY"))
    return false
  }
  return true
}))

nn.MustApplyMiddleware(node.RemoteNodeConnected(func(remoteNode *node.RemoteNode) bool {
  log.Infof("Only lucky remote node can get here :)")
  return true
}))
```

Middleware itself is stateless, but very likely you may need a stateful
middleware for more complex logic. Stateful middleware can be created in a
variety ways without introducing more complex API. One of the simplest ways is
just calling a method in the middleware. As a conceptual example

```go
coolObj := NewCoolObj()
nn.MustApplyMiddleware(node.BytesReceived(func(msg, msgID, srcID []byte, remoteNode *node.RemoteNode) ([]byte, bool) {
  coolObj.DoSomeCoolStuff(msg)
  return msg, true
}))
```

Stateful middleware if very powerful and can give you almost full control of the
system, and it can define its own middleware for other components to use.
Actually the whole overlay network can be implemented as a stateful middleware
that hooked into a few places like `node.LocalNodeWillStart`. We didn't choose
to do it because overlay network is a top-level type in nnet, but there is
nothing that prevent you to do it.

There are lots of middleware types that can be (and should be) used to listen to
and control topology change, message routing and handling, etc. Some of them
provide convenient shortcuts while some provide detailed low level control.
Currently middleware type declaration are distributed in a few places:

* [node/middleware.go](node/middleware.go)
* [overlay/middleware.go](overlay/middleware.go)
* [overlay/chord/middleware.go](overlay/chord/middleware.go)
* [overlay/routing/middleware.go](overlay/routing/middleware.go)

Middleware architecture is very flexible and new type of middleware can be added
easily without breaking existing code. Feel free to [open an
issue](https://github.com/nknorg/nnet/issues/new) if you feel the need for new
middleware type.

A complete middleware example can be found at
[examples/middleware/main.go](examples/middleware/main.go). You can run it
by

```bash
go run $GOPATH/src/github.com/nknorg/nnet/examples/middleware/main.go
```

Also, examples like [examples/message/main.go](examples/message/main.go) and
[examples/efficient-broadcasting/main.go](examples/efficient-broadcasting/main.go)
use middleware to handle and count messages.

### Sending and Receiving Messages

Sending message is top-level function in nnet. You can send arbitrary bytes to
any other node/nodes by calling one of `.SendXXX` methods. Currently there are 3
types of user-defined message classified by routing types:

* Direct: message will be sent to a remote node that has a direct connection with you.
* Relay: message will be routed and delivered to the node with a certain ID, or the node whose ID is closest to the destination ID, typically not directly connected with you. Relay message is routed using DHT topology.
* Broadcast: message will be routed and delivered to every node in the network, not just the nodes you are directly connected to.

The broadcast message has a few subtypes:

* Push message that use simple flooding/gossip protocol and try to send message to every neighbors. Each node will receive the same message C times where C is its neighbor count. **Push message is optimal in terms of latency and robustness and is ideal for small piece of important message** (like votes in consensus).
* Pull message that use push message to send message hash first. Node receiving a message hash and do not have the message itself will pull the message from the neighbor that sent it the hash. Each node will receive the same message hash C times but will only receive the message itself once, with a round trip delay added for each hop. **Pull message is optimal in terms of bandwidth and robustness and is ideal for large piece of important message** (like blocks in blockchain). (to be implemented)
* Tree message that send the message through the spanning tree constructed by the Chord topology. Each node will only receive the same message K times where K is adjustable and can be as small as 1. **Tree message is optimal in terms of both bandwidth and latency but is less robust, and is ideal for small piece of not-that-important information** (like transactions in blockchain).

nnet uses router architecture such that implementing a new routing algorithm is
as simple as implementing a `Router` interface defined in
[routing/routing.go](routing/routing.go) that computes the next hop using
`GetNodeToRoute` method and register the routing type by calling
`localNode.RegisterRoutingType` if needed.

For each routing types, there are 2 sending message APIs: **async** where send
call is immediately returned if send success, and **sync** that will be blocked
and wait for reply or timeout before return. Under the hood, sync mode creates a
reply channel and waits for the message that replies to the original message ID
so you can safely use it while sending/receiving other messages at the same
time.

Sending an async message is straightforward. Let's send a relay message as an
example:

```go
success, err := nn.SendBytesRelayAsync([]byte("Hello world!"), destID)
```

You can choose to send the message in sync way such that the function call will
only return after receiving the reply or timeout:

```go
reply, senderID, err := nn.SendBytesRelaySync([]byte("Hello world!"), destID)
```

To handle received message and send back reply message, we can use the
`node.BytesReceived` middleware together with `SendBytesRelayReply` method.

```go
nn.MustApplyMiddleware(node.BytesReceived(func(msg, msgID, srcID []byte, remoteNode *node.RemoteNode) ([]byte, bool) {
  nn.SendBytesRelayReply(msgID, []byte("Well received!"), srcID) // error is omitted here for simplicity
  return msg, true
}))
```

A complete message example can be found at
[examples/message/main.go](examples/message/main.go). You can run it
by

```bash
go run $GOPATH/src/github.com/nknorg/nnet/examples/message/main.go
```

There is another example
[examples/efficient-broadcasting/main.go](examples/efficient-broadcasting/main.go)
that counts and compares the received message count for push and tree broadcast
message. You will see how tree message can reduce the bandwidth usage by an
order of magnitude. You can run it by

```bash
go run $GOPATH/src/github.com/nknorg/nnet/examples/efficient-broadcasting/main.go
```

### Transport protocol

Transport layer is a separate layer that is part of a node in nnet. Each node
can independently choose what transport protocol it listens to, and is able to
talk to nodes that listen to different transport protocol as long as the
corresponding transport layer is implemented. Each node address starts with the
transport protocol that the node listens to, e.g. `tcp://127.0.0.1:23333`, such
that other nodes know what protocol to use when talking to it.

Currently nnet have 2 transport protocol implemented: TCP and
[KCP]((https://github.com/skywind3000/kcp/blob/master/README.en.md)) (a reliable
low-latency protocol based on UDP). In theory, any other reliable protocol can
be easily integrated by implementing `Dial` and `Listen` interface. Feel free to
[open an issue](https://github.com/nknorg/nnet/issues/new) if you feel the need
for new transport protocol.

Changing transport protocol is as simple as changing the `Transport` value in
config when creating nnet. A complete example can be found at
[examples/mixed-protocol/main.go](examples/mixed-protocol/main.go). You can run
it by

```bash
go run $GOPATH/src/github.com/nknorg/nnet/examples/mixed-protocol/main.go
```

### NAT Traversal

If you are developing an application that is open to public, it is very likely
that you need some sort of NAT traversal since lots of devices are behind NAT at
the moment. NAT traversal can be set up in middleware such as
`overlay.NetworkWillStart` and `node.LocalNodeWillStart`, or anytime before the
network starts by calling the `SetInternalPort` method of the LocalNode. A
complete NAT example that works for UPnP or NAT-PMP enabled routers can be found
at [examples/nat/main.go](examples/nat/main.go). You can run it by

```bash
go run $GOPATH/src/github.com/nknorg/nnet/examples/nat/main.go
```

### Logger

Typically when you use nnet as the network layer of your application, you
probably want it to share the logger with the rest of your application. You can
do it by calling `nnet.SetLogger(logger)` method as long as `logger` implements
the `Logger` interface defined in [log/log.go](log/log.go). If you don't set it,
nnet will use [go-logging](github.com/op/go-logging) by default.

## Benchmark

Throughput is a very important metric of the network stack. There are multiple
potential bottlenecks when we consider throughput, e.g. network i/o, message
serialization/deserialization, architecture and implementation efficiency. To
measure the throughput, we wrote a simple local benchmark, which can be found at
[examples/message-benchmark/main.go](examples/message-benchmark/main.go). You
can run it with default arguments (2 nodes, 1 KB message size) by

```bash
go run $GOPATH/src/github.com/nknorg/nnet/examples/message-benchmark/main.go
```

On a MacBook Pro 2018 this will give you around 75k message/s per node,
far more than enough for most p2p applications.

When message size increase, the bottleneck will becomes bandwidth rather than
message count. We can run the benchmark with 1 MB message by

```bash
go run $GOPATH/src/github.com/nknorg/nnet/examples/message-benchmark/main.go -m 1048576
```

This will give you around 900 MB/s per node on the same computer, again far
more than enough in typical cases.

The same benchmark can be used to see how spanning tree broadcasting can greatly
improve throughput. We first use the standard push method to broadcast in a
16-node network:

```bash
go run $GOPATH/src/github.com/nknorg/nnet/examples/message-benchmark/main.go -n 16 -b push
```

Each node can only receive around 300 unique message/s. The number is so low
because: 1. we are running all 16 nodes on the same laptop with just 6 cores and
more importantly 2. Gossip protocol has a high message redundancy, and in this
particular example each node receives the same message 15 times! If we sum them
up, the total messages being processed are 300*15*15, which is pretty close to
the 75k we got in 2-node benchmark. Of course if the overlay topology is fully
connected, we don't need Gossip protocol at all, but similar inefficiency will
happen once the network size is large enough so that fully connected topology is
unrealistic.

Now let's change to spanning tree broadcasting:

```bash
go run $GOPATH/src/github.com/nknorg/nnet/examples/message-benchmark/main.go -n 16 -b tree
```

Now each node can receive 3k unique messages/s, an order of magnitude boost!
Actually if the benchmark is running on multiple computers with enough network
resources, each node can again receive around 75k messages/s.

## Who is using nnet

* [NKN](https://github.com/nknorg/nkn): a blockchain-powered decentralized data transmission network that uses nnet as network layer.

Welcome to open a [pull request](https://github.com/nknorg/nkn/pulls) if you
want your project using nnet to be listed here. The project needs to be
maintained in order to be listed.

## Contributing

**Can I submit a bug, suggestion or feature request?**

Yes. Please [open an issue](https://github.com/nknorg/nnet/issues/new) for that.

**Can I contribute patches to NKN project?**

Yes, we appreciate your help! To make contributions, please fork the repo, push
your changes to the forked repo with signed-off commits, and open a pull request
here.

Please sign off your commit by adding -s when committing:

```shell
git commit -s
```

This means git will automatically add a line "Signed-off-by: Name <email>" at
the end of each commit, indicating that you wrote the code and have the right to
pass it on as an open source patch.

## Community

* [Discord](https://discord.gg/c7mTynX)
