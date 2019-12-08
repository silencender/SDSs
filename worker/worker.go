package worker

import (
	"log"
	"net"

	. "github.com/silencender/SDSs/utils"
)

type Worker struct {
	port int
}

func NewWorker(port int) *Worker {
	worker := &Worker{
		port: port,
	}
	return worker
}

func (worker *Worker) StartWorker() {
	conn, err := net.Dial("tcp", MasterAddrToW)
	PrintIfErr(err)
	worker_node := &WorkerNode{
		master:     NewNode(conn),
		registered: make(chan *Node),
		unregister: make(chan *Node),
	}
	worker_node.master.Open()
	//之后注册完之后负责关闭连接
	worker_node.register(worker.port)
	log.Println("register done")
	//负责listen
	go worker_node.listen(worker.port)
	//负责register
	go worker_node.run()
}
