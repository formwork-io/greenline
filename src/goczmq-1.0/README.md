# goczmq [![Build Status](https://travis-ci.org/zeromq/goczmq.svg?branch=master)](https://travis-ci.org/zeromq/goczmq) [![Doc Status](https://godoc.org/github.com/zeromq/goczmq?status.png)](https://godoc.org/github.com/zeromq/goczmq)

## Introduction
A golang interface to the [CZMQ v3](http://czmq.zeromq.org) API.

## Update 2016-04-23 - gopkg.in
Support for new thread safe socket types was recently added to CZMQ:

New Socket Type| Old Socket Type
---------------|----------------
ZMQ_RADIO      | ZMQ_PUB
ZMQ_DISH       | ZMQ_SUB
ZMQ_CLIENT     | ZMQ_DEALER
ZMQ_SERVER     | ZMQ_ROUTER
ZMQ_SCATTER    | ZMQ_PUSH
ZMQ_GATHER     | ZMQ_PULL

Thread safe sockets are of immense value as far as Golang goes. Single frame messages
also makes supporting io.Reader / io.Writer a much simpler proposition. However, 
support for these sockets requires building against bleeding edge (as of this date)
czmq / libzmq. There may be API breaking changes as a result as well. In order to 
allow people to continue using GoCZMQ with the current stable releases of libzmq
and czmq, we are moving to support [gopkg.in](http://labix.org/gopkg.in) urls.

### Releases

CZMQ Version | Package Url
-------------|------------
<= 3.02      | [`gopkg.in/zeromq/cmzq.v1`](https://gopkg.in/zeromq/czmq.v1)

### Development Changes

Work for GoCZMQ v2 will take place on a v2.0 branch. The C4.1 development process
will still be used, other than the ammendent that PRs for 2.0 work should be
opened from a fork against branch v2.0 rather than master. I am torn at this additional
complication, but I believe it is the best way to move forward.

## Update 2015-06-02
With the releases of zeromq 4.1 and czmq 3, we are declaring a stable release.
What does "stable release" mean in the world of go getting projects off of git 
master? It means our exported API is finalized and we will, to the absolute
best of our ability, make no breaking changes to it. We may make additive
changes, and we may refactor internals in order to improve performance. 

## Installation

### Building From Source (Linux)

```
wget https://download.libsodium.org/libsodium/releases/libsodium-1.0.10.tar.gz
wget https://download.libsodium.org/libsodium/releases/libsodium-1.0.10.tar.gz.sig
wget https://download.libsodium.org/jedi.gpg.asc
gpg --import jedi.gpg.asc
gpg --verify libsodium-1.0.10.tar.gz.sig libsodium-1.0.10.tar.gz
tar zxvf libsodium-1.0.10.tar.gz
cd libsodium-1.010.
./configure; make check
sudo make install
sudo ldconfig
```

```
wget http://download.zeromq.org/zeromq-4.1.1.tar.gz
tar zxvf zeromq-4.1.1.tar.gz
cd zeromq-4.1.1
./configure --with-libsodium; make; make check
sudo make install
sudo ldconfig
```

```
wget http://download.zeromq.org/czmq-3.0.1.tar.gz
tar zxvf czmq-3.0.1.tar.gz
cd czmq-3.0.1
./configure; make check
sudo make install
sudo ldconfig
```

```
go get github.com/zeromq/goczmq
```

### Installing on OSX

```
brew install zmq czmq libsodium
```

```
go get github.com/zeromq/goczmq
```

## Usage
### Direct CZMQ Sock API
#### Example
```go
package main

import (
	"log"

	"github.com/zeromq/goczmq"
)

func main() {
	// Create a router socket and bind it to port 5555.
	router, err := goczmq.NewRouter("tcp://*:5555")
	if err != nil {
		log.Fatal(err)
	}
	defer router.Destroy()

	log.Println("router created and bound")

	// Create a dealer socket and connect it to the router.
	dealer, err := goczmq.NewDealer("tcp://127.0.0.1:5555")
	if err != nil {
		log.Fatal(err)
	}
	defer dealer.Destroy()

	log.Println("dealer created and connected")

	// Send a 'Hello' message from the dealer to the router.
	// Here we send it as a frame ([]byte), with a FlagNone
	// flag to indicate there are no more frames following.
	err = dealer.SendFrame([]byte("Hello"), goczmq.FlagNone)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("dealer sent 'Hello'")

	// Receve the message. Here we call RecvMessage, which
	// will return the message as a slice of frames ([][]byte).
	// Since this is a router socket that support async
	// request / reply, the first frame of the message will
	// be the routing frame.
	request, err := router.RecvMessage()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("router received '%s' from '%v'", request[1], request[0])

	// Send a reply. First we send the routing frame, which
	// lets the dealer know which client to send the message.
	// The FlagMore flag tells the router there will be more
	// frames in this message.
	err = router.SendFrame(request[0], goczmq.FlagMore)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("router sent 'World'")

	// Next send the reply. The FlagNone flag tells the router
	// that this is the last frame of the message.
	err = router.SendFrame([]byte("World"), goczmq.FlagNone)
	if err != nil {
		log.Fatal(err)
	}

	// Receive the reply.
	reply, err := dealer.RecvMessage()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("dealer received '%s'", string(reply[0]))
}
```
#### Output
```
2015/05/26 21:52:52 router created and bound
2015/05/26 21:52:52 dealer created and connected
2015/05/26 21:52:52 dealer sent 'Hello'
2015/05/26 21:52:52 router received 'Hello' from '[0 103 84 189 175]'
2015/05/26 21:52:52 router sent 'World'
2015/05/26 21:52:52 dealer received 'World'
```
### io.ReadWriter support
#### Example
```go
package main

import (
	"log"

	"github.com/zeromq/goczmq"
)

func main() {
	// Create a router socket and bind it to port 5555.
	router, err := goczmq.NewRouter("tcp://*:5555")
	if err != nil {
		log.Fatal(err)
	}
	defer router.Destroy()

	log.Println("router created and bound")

	// Create a dealer socket and connect it to the router.
	dealer, err := goczmq.NewDealer("tcp://127.0.0.1:5555")
	if err != nil {
		log.Fatal(err)
	}
	defer dealer.Destroy()

	log.Println("dealer created and connected")

	// Send a 'Hello' message from the dealer to the router,
	// using the io.Write interface
	n, err := dealer.Write([]byte("Hello"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("dealer sent %d byte message 'Hello'\n", n)

	// Make a byte slice and pass it to the router
	// Read interface. When using the ReadWriter
	// interface with a router socket, the router
	// caches the routing frames internally in a
	// FIFO and uses them transparently when
	// sending replies.
	buf := make([]byte, 16386)

	n, err = router.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("router received '%s'\n", buf[:n])

	// Send a reply.
	n, err = router.Write([]byte("World"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("router sent %d byte message 'World'\n", n)

	// Receive the reply, reusing the previous buffer.
	n, err = dealer.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("dealer received '%s'", string(buf[:n]))
}
```
#### Output
```
2015/05/26 21:54:10 router created and bound
2015/05/26 21:54:10 dealer created and connected
2015/05/26 21:54:10 dealer sent 5 byte message 'Hello'
2015/05/26 21:54:10 router received 'Hello'
2015/05/26 21:54:10 router sent 5 byte message 'World'
2015/05/26 21:54:10 dealer received 'World'
```
### Thread safe channel interface
#### Example
```go
package main

import (
	"log"

	"github.com/zeromq/goczmq"
)

func main() {
	// Create a router channeler and bind it to port 5555.
	// A channeler provides a thread safe channel interface
	// to a *Sock
	router := goczmq.NewRouterChanneler("tcp://*:5555")
	defer router.Destroy()

	log.Println("router created and bound")

	// Create a dealer channeler and connect it to the router.
	dealer := goczmq.NewDealerChanneler("tcp://127.0.0.1:5555")
	defer dealer.Destroy()

	log.Println("dealer created and connected")

	// Send a 'Hello' message from the dealer to the router.
	dealer.SendChan <- [][]byte{[]byte("Hello")}
	log.Println("dealer sent 'Hello'")

	// Receve the message as a [][]byte. Since this is
	// a router, the first frame of the message wil
	// be the routing frame.
	request := <-router.RecvChan
	log.Printf("router received '%s' from '%v'", request[1], request[0])

	// Send a reply. First we send the routing frame, which
	// lets the dealer know which client to send the message.
	router.SendChan <- [][]byte{request[0], []byte("World")}
	log.Printf("router sent 'World'")

	// Receive the reply.
	reply := <-dealer.RecvChan
	log.Printf("dealer received '%s'", string(reply[0]))
}
```
#### Output
```
2015/05/26 21:56:43 router created and bound
2015/05/26 21:56:43 dealer created and connected
2015/05/26 21:56:43 dealer sent 'Hello'
2015/05/26 21:56:43 received 'Hello' from '[0 12 109 153 35]'
2015/05/26 21:56:43 router sent 'World'
2015/05/26 21:56:43 dealer received 'World'
```
## GoDoc
[godoc](https://godoc.org/github.com/zeromq/goczmq)

## See Also
* [Peter Kleiweg's](https://github.com/pebbe) [zmq4](https://github.com/pebbe/zmq4) bindings

## License
This project uses the MPL v2 license, see LICENSE
