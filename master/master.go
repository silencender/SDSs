package master

import (
	"container/list"
	"log"

	. "github.com/silencender/SDSs/utils"
)

func StartMaster() {
	log.Println("Master running")
	wm := WorkerManager{
		workers:    list.New(),
		pworker:    nil,
		register:   make(chan string),
		unregister: make(chan *Node),
	}
	cm := ClientManager{
		wm:         &wm,
		register:   make(chan *Node),
		unregister: make(chan *Node),
	}
	go cm.run()
	go cm.listen(MasterAddrToC)
	go wm.run()
	go wm.listen(MasterAddrToW)
}
