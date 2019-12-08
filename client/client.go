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

func NewClient(seq int) (*Client){
    log.Println("start worker ",seq)
    client := &Client{
        seq : seq,
        repeatTime : 100,
    }
    client.StartClient()
    return client
}

func (client *Client)StartClient(){
    conn, err := net.Dial("tcp", MasterAddrToC)
    PrintIfErr(err)
    master_node := NewNode(conn)
    master_node.Open()
    cn := &ClientNode{
		Master : *master_node,
		WorkerList: make(chan *Node),
        QueryList: make(chan []byte),
        Pool: WorkerPool{
			workers:    make(map[string]*Node),
			register:   make(chan *Node),
			unregister: make(chan []byte),
		},
	}
    //返回client结构体
    go cn.generate(client.repeatTime)
    //go cn.query(client.repeatTime)
    //receive 和send如果不在进程池中才建立
}
