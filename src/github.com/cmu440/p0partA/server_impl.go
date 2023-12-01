// Implementation of a KeyValueServer. Students should write their code in this file.

package p0partA

import (
	"fmt"
	"log"
	"net"

	"github.com/cmu440/p0partA/kvstore"
)

type keyValueServer struct {
	// TODO: implement this!

	kvs kvstore.KVStore

	port int

	started       bool
	closed        bool
	count_active  int
	count_dropped int
}

// New creates and returns (but does not start) a new KeyValueServer.
func New(store kvstore.KVStore) KeyValueServer {
	// TODO: implement this!
	return &keyValueServer{kvs: store, started: false, closed: false, count_active: 0, count_dropped: 0}

}

func (kvs *keyValueServer) Start(port int) error {
	// TODO: implement this!
	if kvs.closed {
		return fmt.Errorf("server has been closed")
	}
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	kvs.started = true

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go kvs.HandleConn(conn)
		}
	}()

	return nil
}

func (kvs *keyValueServer) Close() {
	kvs.closed = true
}

func (kvs *keyValueServer) CountActive() int {
	return kvs.count_active
}

func (kvs *keyValueServer) CountDropped() int {
	return kvs.count_dropped
}

// TODO: add additional methods/functions below!

func (kvs *keyValueServer) HandleConn(conn net.Conn) {

}
