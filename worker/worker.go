package worker

import (
    . "github.com/silencender/SDSs/utils"
    "log"
    "net"
)

type Worker struct {
	port string
}

func NewWorker(port string) *Worker {
	worker := Worker{port: port}
	return &worker
}

func (worker *Worker) StartWorker() {
    log.Println("Worker running in ",worker.port)
    conn,err := net.Dial("tcp",MasterAddrToW)
    PrintIfErr(err)
    worker_node := &WorkerNode{
        master : NewNode(conn),
    }
    worker_node.master.Open()
    //之后注册完之后负责关闭连接
    worker_node.register("18888")
    //负责listen
    //go worker_node.listen(worker.addr)
    //go worker_node.run()
}
