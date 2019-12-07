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

func NewWorker(addr int) *Worker {
	worker := &Worker{
        port: addr,
    }
	return worker
}

func (worker *Worker) StartWorker() {
    local_addr := &net.TCPAddr{Port:worker.port}
    worker_node := &WorkerNode{}
    go worker_node.listen(worker.port)
    d := net.Dialer{LocalAddr: local_addr}
    conn,err := d.Dial("tcp",MasterAddrToW)
    PrintIfErr(err)
    worker_node.master = NewNode(conn)
    //worker_node := &WorkerNode{
    //    master : NewNode(conn),
    //}
    worker_node.master.Open()
    //之后注册完之后负责关闭连接
    worker_node.register()
    log.Println("register done")
    //负责listen
}
