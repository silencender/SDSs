package master

import (
	"container/list"
	"log"

	. "github.com/silencender/SDSs/utils"
)

type Master struct {
	clientListenAddr string
	workerListenAddr string
}

func NewMaster(cla, wla string) *Master {
	return &Master{
		clientListenAddr: cla,
		workerListenAddr: wla,
	}
}

func (master *Master) StartMaster() {
	log.Println("Master running")
	wm := WorkerManager{
		workers:    list.New(),
		pworker:    nil,
		register:   make(chan *Node),
		unregister: make(chan *Node),
	}
	cm := ClientManager{
		wm:         &wm,
		register:   make(chan *Node),
		unregister: make(chan *Node),
	}
	go cm.run()
	go cm.listen(master.clientListenAddr)
	go wm.run()
	go wm.listen(master.workerListenAddr)
}
