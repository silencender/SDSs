package client

import (
	"log"
	"net"

	. "github.com/silencender/SDSs/utils"
)

type Client struct {
	repeatTime int
	seq        int
	masterAddr string
}

func NewClient(seq, repeats int, ma string) *Client {
	log.Println("start worker ", seq)
	client := &Client{
		seq:        seq,
		repeatTime: repeats,
		masterAddr: ma,
	}
	return client
}

func (client *Client) StartClient() {
	conn, err := net.Dial("tcp", client.masterAddr)
	PrintIfErr(err)
	master_node := NewNode(conn)
	master_node.Open()
	cn := &ClientNode{
		Master:     master_node,
		WorkerList: make(chan *Node),
		QueryList:  make(chan []byte),
		Pool: WorkerPool{
			workers: make(map[string]*Node),
		},
		register:   make(chan *Node),
		unregister: make(chan *Node),
	}
	go cn.run()
	go cn.query(client.repeatTime)
	//建立与主进程通信的receive和send
	go cn.receive(cn.Master)
	go cn.handle(cn.Master)
	go cn.send(cn.Master)
	//receive 和send如果不在进程池中才建立
}
