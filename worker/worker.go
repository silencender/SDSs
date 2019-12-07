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
	worker := &Worker{
        addr: addr,
    }
	return worker
}

func (worker *Worker) StartWorker() {
    conn,err := net.Dial("tcp",MasterAddrToW)
    PrintIfErr(err)
    worker_node := &WorkerNode{
        master : NewNode(conn),
    }
    worker_node.master.Open()
    //之后注册完之后负责关闭连接
    worker_node.register()
    log.Println("register done")
    //负责listen
    go worker_node.listen()
}
