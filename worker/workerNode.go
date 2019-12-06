package worker

import (
	. "github.com/silencender/SDSs/utils"
	pb "github.com/silencender/SDSs/protos"
    "github.com/golang/protobuf/proto"

    "time"
    "log"
    "net"
)

type WorkerNode struct {
	master *Node
}

func (wn *WorkerNode) register(port string) {
	log.Println("register to worker for port ",port)
    registReq := &pb.Message{
		MsgType:pb.Message_REGISTER_REQ,
		Seq: int32(time.Now().Unix()),
		Socket: port,
	}
    registReqData,err := proto.Marshal(registReq)
    PrintIfErr(err)
	wn.master.Socket.Write([]byte(registReqData))
	//事实上不用接到master的反馈也行，虽然定义了

}

func (wn *WorkerNode) receive(client *Node) {
    message := make([]byte,BufSize)
    for {
        length,err :=client.Socket.Read(message)
        PrintIfErr(err)
        if length >0 {
            log.Println("received ",length," bytes from ",client.Socket.RemoteAddr)
        }
    }
}

func (wn *WorkerNode) handle(client *Node) {

}

func (wn *WorkerNode) send(client *Node) {

}

func (wn *WorkerNode) listen(port string) {
    listener,err :=net.Listen("tcp","localhost:"+port)
    PrintIfErr(err)
    for{
        conn,err := listener.Accept()
        PrintIfErr(err)
        log.Println("listened a connect from ",conn.RemoteAddr)
        client := NewNode(conn)
        go wn.receive(client)
        //go wn.handle(client)
        //go wn.send(client)
    }
}

func (wn *WorkerNode) run() {

}
