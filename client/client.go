package client

import (
	. "github.com/silencender/SDSs/utils"
)

type Client struct {
	master     Node
	queryList  chan []byte
	workerList chan *Node
}

type WorkerPool struct {
	workers    map[string]*Node
	register   chan *Node
	unregister chan *Node
}

func (client *Client) register() {

}

func (client *Client) query() {

}

func (client *Client) receive(worker *Node) {

}

func (client *Client) send(worker *Node) {

}

func run() {

}

func startClient() {

}
