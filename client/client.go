package client

import (
    . "github.com/silencender/SDSs/utils"
    "log"
    "net"
)

type Client struct{
    repeatTime int
    seq int
}

func NewClient(seq,repeats int) (*Client){
    log.Println("start worker ",seq)
    client := &Client{
        seq : seq,
        repeatTime : repeats,
    }
    return client
}

func (client *Client)StartClient(){
    conn, err := net.Dial("tcp", MasterAddrToC)
    PrintIfErr(err)
    master_node := NewNode(conn)
    master_node.Open()
    cn := &ClientNode{
		Master : master_node,
		WorkerList: make(chan *Node),
        QueryList: make(chan []byte),
        Pool: WorkerPool{
			workers:    make(map[string]*Node),
		},
		register:   make(chan *Node),
		unregister: make(chan *Node),
	}
    go cn.run(client.repeatTime)
    go cn.query(client.repeatTime)
    //建立与主进程通信的receive和send
    go cn.receive(cn.Master)
    go cn.handle(cn.Master)
    go cn.send(cn.Master)
    go cn.registerManager()
    //receive 和send如果不在进程池中才建立
}
