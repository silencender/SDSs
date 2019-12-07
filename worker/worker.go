package worker

import (
    . "github.com/silencender/SDSs/utils"
    "log"
    "net"
    //"strings"
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
    conn,err := net.Dial("tcp",MasterAddrToW)
    PrintIfErr(err)
    worker_node := &WorkerNode{
        master : NewNode(conn),
    }
    worker_node.master.Open()
    //之后注册完之后负责关闭连接
    worker_node.register(worker.port)
    log.Println("register done")
    //负责listen
    go worker_node.listen(worker.port)
}
