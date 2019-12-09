package worker

import (
	"log"
	"net"

	. "github.com/silencender/SDSs/utils"
)

type Worker struct {
	addr       string
	masterAddr string
}

func NewWorker(addr, ma string) *Worker {
	worker := &Worker{
		addr:       addr,
		masterAddr: ma,
	}
	return worker
}

func (worker *Worker) StartWorker() {
	conn, err := net.Dial("tcp", worker.masterAddr)
	PrintIfErr(err)
	worker_node := &WorkerNode{
		master:     NewNode(conn),
		registered: make(chan *Node),
		unregister: make(chan *Node),
	}
	worker_node.master.Open()
	//之后注册完之后负责关闭连接
	worker_node.register(worker.addr)
	log.Println("Register done")
	//负责listen
	go worker_node.listen(worker.addr)
	//负责register
	go worker_node.run()
}
