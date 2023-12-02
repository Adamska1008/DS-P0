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

	started       bool
	closed        bool
	count_active  int
	count_dropped int
}

// New creates and returns (but does not start) a new KeyValueServer.
func New(store kvstore.KVStore) KeyValueServer {
	// TODO: implement this!
	return &keyValueServer{store: store, started: false, closed: false, count_active: 0, count_dropped: 0}

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
			kvs.count_active++
			if err != nil {
				log.Fatal("line50 ", err)
			}
			go kvs.HandleConn(conn)
			kvs.count_dropped++
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
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		message, err := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		fmt.Println("Read Message: ", message)
		if err != nil {
			conn.Close()
			if err == io.EOF {
				break
			}
			log.Fatal("line 86 ", err)
			return
		}
		args := strings.Split(message, ":")
		switch args[0] {
		case "Get": // Get:[key]
			value := kvs.store.Get(args[1])
			fmt.Printf("Get key %v\n", args[1])
			for i := range value {
				_, err := writer.Write(value[i])
				if err != nil {
					log.Fatal("line 97 ", err)
				}
				writer.WriteByte('\n')
			}
			writer.Flush()
		case "Delete": // Delete:[key]
			// fmt.Printf("Delete key %v\n", args[1])
			kvs.store.Delete(args[1])
		case "Update": // Update:[key]:[oldValCleaned]:[newVal]
			// fmt.Printf("Update %v %v %v\n", args[1], args[2], args[3])
			kvs.store.Update(args[1], []byte(args[2]), []byte(args[3]))
		case "Put": // Put:[key]:[val]
			args := strings.Split(message, ":")
			// fmt.Printf("Put %v %v\n", args[1], args[2])
			kvs.store.Put(args[1], []byte(args[2]))
		default:
			log.Fatal("line 112 no command: ", args[0])
		}
	}
}
