package master

import (
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	pb "github.com/silencender/SDSs/protos"
	. "github.com/silencender/SDSs/utils"
)

type ClientManager struct {
	wm         *WorkerManager
	register   chan *Node
	unregister chan *Node
}

func (cm *ClientManager) receive(client *Node) {
	message := make([]byte, BufSize)
	for {
		length, err := client.Socket.Read(message)
		if err != nil {
			cm.unregister <- client
			close(client.ReqData)
			break
		}
		if length > 0 {
			client.ReqData <- message
		}
	}

}

func (cm *ClientManager) handle(client *Node) {
	for {
		select {
		case req, ok := <-client.ReqData:
			if !ok {
				close(client.ResData)
				return
			}
			message := &pb.Message{}
			err := proto.Unmarshal(req, message)
			PrintIfErr(err)
			res := &pb.Message{
				Seq: message.GetSeq(),
			}
			switch message.MsgType {
			case pb.Message_REGISTER_REQ:
				res.MsgType = pb.Message_REGISTER_RES
			case pb.Message_QUERY_REQ:
				res.MsgType = pb.Message_QUERY_RES
				res.Socket = cm.wm.SelectWorker().Info.String()
			}
			data, err := proto.Marshal(res)
			PrintIfErr(err)
			client.ResData <- data
		}
	}
}

func (cm *ClientManager) send(client *Node) {
	for {
		select {
		case message, ok := <-client.ResData:
			if !ok {
				return
			}
			client.Socket.Write(message)
		}
	}
}

func (cm *ClientManager) listen(addr string) {
	listener, err := net.Listen("tcp", addr)
	PrintIfErr(err)
	for {
		conn, err := listener.Accept()
		PrintIfErr(err)
		client := NewNode(conn)
		cm.register <- client
		go cm.receive(client)
		go cm.handle(client)
		go cm.send(client)
	}
}

func (cm *ClientManager) run() {
	for {
		select {
		case conn := <-cm.register:
			conn.Open()
			log.Printf("Client %s registered\n", conn.Info.String())
		case conn := <-cm.unregister:
			conn.Close()
			log.Printf("Client %s unregistered\n", conn.Info.String())
		}
	}
}
