package client

import (
    . "SDSs/utils"
	"SDSs/protos"
	"github.com/golang/protobuf/proto"
    "fmt"
    "net"
    "time"
)

type WorkerPool struct {
	workers    map[string]*Node
	register   chan *Node
	unregister chan *Node
}

type Client struct {
	Master     Node
	Pool    WorkerPool
    QueryList  chan []byte
	WorkerList chan *Node
}


func (client *Client) register() {

}

func (client *Client) query() {
	queryReq := &protos.Message{
		MsgType:protos.Message_QUERY_REQ,
		Seq: int32(time.Now().Unix()),
    } 
	queryReqData,err := proto.Marshal(queryReq)
	if err != nil {
		fmt.Println(err)
	}
	client.Master.Socket.Write([]byte(queryReqData))
}

func (client *Client) receive(worker *Node) {

}

func (client *Client) send(worker *Node) {

}

//数据类型说明
//calcType:'f,i,l,d'
//calcOp1\calcOp2:对应的运算数
func (client *Client) Run(calcType,calcOp1,calcOp2 string) {
    //calcString := calcType + ":" + calcOp1 + ":" + calcOp2
    client.query()
}

func StartClient() Client {
    //建立与master的连接
    conn, err := net.Dial("tcp", MasterAddr)
	if err != nil {
		fmt.Println("net.Dial err = ", err)
        //初始化空类client
        client := Client{}
		return client 
    }
    client := Client{
		Master: Node{
		    Socket: conn,
		    Ok:     true,
		    Info:   conn.RemoteAddr(),
		    Data:   make(chan []byte),
		},
		WorkerList: make(chan *Node),
		Pool: WorkerPool{
			workers:    make(map[string]*Node),
			register:   make(chan *Node),
			unregister: make(chan *Node),
		},
	}
    fmt.Println("done")
    //返回client结构体
    return client
}
