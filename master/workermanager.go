package master

import (
	"container/list"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	pb "github.com/silencender/SDSs/protos"
	. "github.com/silencender/SDSs/utils"
)

type WorkerManager struct {
	workers    *list.List
	pworker    *list.Element
	register   chan *Node
	unregister chan *Node
}

func (wm *WorkerManager) receive(worker *Node) {
	message := make([]byte, BufSize)
	for {
		length, err := worker.Socket.Read(message)
		if err != nil {
			wm.unregister <- worker
			close(worker.ReqData)
			break
		}
		if length > 0 {
			worker.ReqData <- message
		}
	}
}

func (wm *WorkerManager) handle(worker *Node) {
	for {
		select {
		case req, ok := <-worker.ReqData:
			if !ok {
				close(worker.ResData)
				return
			}
			message := &pb.Message{}
			err := proto.Unmarshal(req, message)
			PrintIfErr(err)
			res := &pb.Message{
				Seq: message.Seq(),
			}
			switch message.MsgType {
			case pb.Message_REGISTER_REQ:
				res.MsgType = pb.Message_REGISTER_RES
			case pb.Message_HEARTBEAT_REQ:
				res.MsgType = pb.Message_HEARTBEAT_RES
			}
			data, err := proto.Marshal(res)
			PrintIfErr(err)
			worker.ResData <- data
		}
	}
}

func (wm *WorkerManager) send(worker *Node) {
	for {
		select {
		case message, ok := <-worker.ResData:
			if !ok {
				return
			}
			worker.Socket.Write(message)
		}
	}
}

func (wm *WorkerManager) listen(addr string) {
	listener, err := net.Listen("tcp", addr)
	PrintIfErr(err)
	for {
		conn, err := listener.Accept()
		PrintIfErr(err)
		worker := NewNode(conn)
		wm.register <- worker
		go wm.receive(worker)
		go wm.handle(worker)
		go wm.send(worker)
	}
}

func (wm *WorkerManager) run() {
	for {
		select {
		case conn := <-wm.register:
			conn.Open()
			wm.workers.PushBack(conn)
			log.Printf("Worker %s registered\n", conn.Info.String())
		case conn := <-wm.unregister:
			conn.Close()
			RemoveListItem(wm.workers, conn)
			log.Printf("Worker %s unregistered\n", conn.Info.String())
		}
	}
}

func (wm *WorkerManager) SelectWorker() *Node {
	if wm.pworker == nil {
		wm.pworker = wm.workers.Front()
	}
	var worker *Node
	for {
		if wm.pworker == wm.workers.Back() {
			wm.pworker = wm.workers.Front()
		} else {
			wm.pworker = wm.pworker.Next()
		}
		worker = wm.pworker.Value.(*Node)
		if worker.Ok {
			return worker
		}
	}
}
