package client

import (
	"fmt"

	. "github.com/silencender/SDSs/utils"
)

type ClientNode struct {
	master     *Node
	workerList chan *Node
	workerpool *WorkerPool
}

type WorkerPool struct {
	workers    map[string]*Node
	register   chan *Node
	unregister chan *Node
}

func (cn *ClientNode) register() {
	cn.master.Data <- []byte("ClientNode hello\n")
}

func (cn *ClientNode) query() {

}

func (cn *ClientNode) receive(worker *Node) {
	for {
		message := make([]byte, 1024)
		length, err := cn.master.Socket.Read(message)
		if err != nil {
			break
		}
		if length > 0 {
			fmt.Printf("Received from master %s: %s", cn.master.Info.String(), string(message))
		}
	}
}

func (cn *ClientNode) send(worker *Node) {
	for {
		select {
		case message, ok := <-cn.master.Data:
			if !ok {
				return
			}
			cn.master.Socket.Write(message)
		}
	}
}

func (cn *ClientNode) run() {

}
