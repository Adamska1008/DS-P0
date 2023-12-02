// Implementation of a KeyValueServer. Students should write their code in this file.

package p0partA

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/cmu440/p0partA/kvstore"
)

type keyValueServer struct {
	store kvstore.KVStore

	started bool
	closed  bool

	addActive      chan bool
	addDropped     chan bool
	activeRequest  chan bool
	activeResult   chan int
	droppedRequest chan bool
	droppedResult  chan int
	countActive    int
	countDropped   int
}

// New creates and returns (but does not start) a new KeyValueServer.
func New(store kvstore.KVStore) KeyValueServer {
	// TODO: implement this!
	return &keyValueServer{
		store:          store,
		started:        false,
		closed:         false,
		addActive:      make(chan bool),
		addDropped:     make(chan bool),
		activeRequest:  make(chan bool),
		activeResult:   make(chan int),
		droppedRequest: make(chan bool),
		droppedResult:  make(chan int),
		countActive:    0,
		countDropped:   0,
	}

}

func (kvs *keyValueServer) Start(port int) error {
	// TODO: implement this!
	if kvs.closed {
		return errors.New("server has been closed")
	}
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("line41 ", err)
	}
	kvs.started = true

	go func() {
		for {
			conn, err := listener.Accept()
			kvs.AddActive()
			if err != nil {
				log.Fatal("line50 ", err)
			}
			kvs.HandleConn(conn)
		}
	}()

	go kvs.chanRountine()

	return nil
}

func (kvs *keyValueServer) chanRountine() { // select-style channel handler for active and dropped
	for {
		select {
		case <-kvs.addActive:
			kvs.countActive++
		case <-kvs.activeRequest:
			kvs.activeResult <- kvs.countActive
		case <-kvs.addDropped:
			kvs.countDropped++
		case <-kvs.droppedRequest:
			kvs.droppedResult <- kvs.countDropped
		}
	}
}

func (kvs *keyValueServer) Close() {
	kvs.closed = true
}

func (kvs *keyValueServer) AddActive() {
	kvs.addActive <- true
}

func (kvs *keyValueServer) AddDropped() {
	kvs.addDropped <- true
}

func (kvs *keyValueServer) CountActive() int {
	kvs.activeRequest <- true
	return <-kvs.activeResult
}

func (kvs *keyValueServer) CountDropped() int {
	kvs.droppedRequest <- true
	return <-kvs.droppedResult
}

// TODO: add additional methods/functions below!

func (kvs *keyValueServer) HandleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		message, err := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		fmt.Println("Read Message: ", message)
		if err != nil {
			conn.Close()
			if err == io.EOF {
				kvs.AddDropped()
				return
			}
			log.Fatal("line 86 ", err)
			return
		}
		args := strings.Split(message, ":")
		switch args[0] {
		case "Get": // Get:[key]
			value := kvs.store.Get(args[1])
			for i := range value {
				_, err := writer.Write(value[i])
				if err != nil {
					log.Fatal("line 96 ", err)
				}
				writer.WriteByte('\n')
			}
			writer.Flush()
		case "Delete": // Delete:[key]
			kvs.store.Delete(args[1])
		case "Update": // Update:[key]:[oldValCleaned]:[newVal]
			fmt.Printf("Update %v %v %v\n", args[1], args[2], args[3])
			kvs.store.Update(args[1], []byte(args[2]), []byte(args[3]))
			kvs.Check(args[1])
		case "Put": // Put:[key]:[val]
			kvs.store.Put(args[1], []byte(args[2]))
		default:
			log.Fatal("line 110 no command: ", args[0])
		}
	}
}

func (kvs *keyValueServer) Check(key string) {
	value := kvs.store.Get(key)
	for _, v := range value {
		fmt.Println(string(v))
	}
}
