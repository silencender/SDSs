package worker

import (
    . "github.com/silencender/SDSs/utils"
    "log"
    "net"
    "strings"
)

type Worker struct {
	addr string
}

func NewWorker() *Worker {
	worker := Worker{addr: ""}
	return &worker
}

func (worker *Worker) StartWorker() {
    conn,err := net.Dial("tcp",MasterAddrToW)
    worker.addr = conn.LocalAddr().String()
    log.Println("Worker running in ",worker.addr)
    port := strings.Split(worker.addr,":")[1]
    PrintIfErr(err)
    worker_node := &WorkerNode{
        master : NewNode(conn),
    }
    worker_node.master.Open()
    //之后注册完之后负责关闭连接
    worker_node.register(port)
    //负责listen
    go worker_node.listen(worker.addr)
    //负责处理注册
    //没有实现，还是没想懂run的必要性
    //go worker_node.run()
}
