package master

import (
	"container/list"

	. "github.com/silencender/SDSs/utils"
)

type WorkerManager struct {
	workers    list.List
	pworker    *list.Element
	register   chan *Node
	unregister chan *Node
}

func (wm *WorkerManager) receive(worker *Node) {

}

func (wm *WorkerManager) send(worker *Node) {

}

func (wm *WorkerManager) run() {

}

func (wm *WorkerManager) selectWorker() *Node {
	return nil
}
