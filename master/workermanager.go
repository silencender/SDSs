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
	register   chan string
	unregister chan *Node
}

func (wm *WorkerManager) receive(worker *Node) {
	message := make([]byte, BufSize)
    log.Println("hhh,I received it")
	length, err := worker.Socket.Read(message)
	if err != nil {
			wm.unregister <- worker
	        return
    }
	if length > 0 {
            worker.ReqData <- message
	}
}

func (wm *WorkerManager) handle(worker *Node) {
	for {
		select {
		case req, ok := <-worker.ReqData:
			if !ok {
				return
			}
            message := &pb.Message{}
			err := proto.Unmarshal(req, message)
			PrintIfErr(err)
			//res := &pb.Message{
			//	Seq: message.Seq,
			//}
			switch message.MsgType {
			case pb.Message_REGISTER_REQ:
                wm.register <- message.Socket
                //res.MsgType = pb.Message_REGISTER_RES
                //case pb.Message_HEARTBEAT_REQ:
				//res.MsgType = pb.Message_HEARTBEAT_RES
			}
			//data, err := proto.Marshal(res)
			PrintIfErr(err)
			//worker.ResData <- data
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
		if err != nil {
			log.Println(err.Error())
		}
		worker := NewNode(conn)
		go wm.handle(worker)
        go wm.receive(worker)
		go wm.send(worker)
	}
}

func (wm *WorkerManager) run() {
	for {
		select {
		case udpAddr := <-wm.register:
			wm.workers.PushBack(udpAddr)
			log.Printf("Worker %s registered\n", udpAddr)
		case conn := <-wm.unregister:
			conn.Close()
			RemoveListItem(wm.workers, conn)
			log.Printf("Worker %s unregistered\n", conn.Info.String())
		}
	}
}

func (wm *WorkerManager) SelectWorker() string {
	if wm.pworker == nil {
		wm.pworker = wm.workers.Front()
	}
	for {
		if wm.pworker == wm.workers.Back() {
			wm.pworker = wm.workers.Front()
		} else {
			wm.pworker = wm.pworker.Next()
		}
        worker := wm.pworker.Value.(string)
        log.Println("I will give u worker: ",worker)
        return worker
	}
}
