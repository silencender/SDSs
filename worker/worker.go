package worker

import (
    . "github.com/silencender/SDSs/utils"
    "log"
    "net"
    //"strings"
)

type Worker struct {
	addr string
}

func NewWorker(addr string) *Worker {
	worker := Worker{addr: addr}
	return &worker
}

func (worker *Worker) StartWorker() {
    conn,err := net.Dial("tcp",MasterAddrToW)
    PrintIfErr(err)
    log.Println("Worker running in ",worker.addr)
    worker_node := &WorkerNode{
        master : NewNode(conn),
    }
    worker_node.master.Open()
    //之后注册完之后负责关闭连接
    worker_node.register(worker.addr)
    //负责listen
    go worker_node.receive(worker.addr)
    //负责返还
    //go worker_node.send()
    //没有实现，还是没想懂run的必要性
    //go worker_node.run()
}
