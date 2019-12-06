package worker

import (
    . "github.com/silencender/SDSs/utils"
    "log"
    "net"
)

type Worker struct {
	addr string
}

func NewWorker(addr string) *Worker {
	worker := Worker{addr: addr}
	return &worker
}

func (worker *Worker) StartWorker(addr1 string) {
    log.Println("Worker running in ",addr1)
    conn,err := net.Dial("tcp",MasterAddrToW)
    PrintIfErr(err)
    worker_node := &WorkerNode{
        master : NewNode(conn),
    }
    worker_node.master.Open()
    go worker_node.receive()
    go worker_node.send()
    go worker_node.run()
}
