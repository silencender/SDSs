package client

import (
	"fmt"
	"net"

	. "github.com/silencender/SDSs/utils"
)

type Client struct {
	master     *Node
	workerList chan *Node
	workerpool *WorkerPool
}

type WorkerPool struct {
	workers    map[string]*Node
	register   chan *Node
	unregister chan *Node
}

func (client *Client) register() {
	client.master.Data <- []byte("Client hello\n")
}

func (client *Client) query() {

}

func (client *Client) receive(worker *Node) {
	for {
		message := make([]byte, 1024)
		length, err := client.master.Socket.Read(message)
		if err != nil {
			break
		}
		if length > 0 {
			fmt.Printf("Received from master %s: %s", client.master.Info.String(), string(message))
		}
	}
}

func (client *Client) send(worker *Node) {
	for {
		select {
		case message, ok := <-client.master.Data:
			if !ok {
				return
			}
			client.master.Socket.Write(message)
		}
	}
}

func (client *Client) run() {

}

func StartClient() {
	fmt.Println("Client running")
	conn, err := net.Dial("tcp", MasterAddrToC)
	if err != nil {
		fmt.Println(err)
	}
	client := Client{
		master: &Node{
			Socket: conn,
			Ok:     false,
			Info:   conn.RemoteAddr(),
			Data:   make(chan []byte),
		},
		workerList: make(chan *Node),
		workerpool: &WorkerPool{
			workers:    make(map[string]*Node),
			register:   make(chan *Node),
			unregister: make(chan *Node),
		},
	}
	go client.run()
	go client.receive(nil)
	go client.send(nil)
	client.register()
}
