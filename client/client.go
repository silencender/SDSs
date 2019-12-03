package client

import (
	"fmt"
	"net"

	. "github.com/silencender/SDSs/utils"
)

type Client struct {
}

func NewClient() *Client {
	client := Client{}
	return &client
}

func (client *Client) StartClient() {
	fmt.Println("Client running")
	conn, err := net.Dial("tcp", MasterAddrToC)
	if err != nil {
		fmt.Println(err)
	}
	cn := ClientNode{
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
	go cn.run()
	go cn.receive(nil)
	go cn.send(nil)
	cn.register()
}
