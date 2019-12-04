package master

import (
	//"container/list"
	"fmt"

	. "SDSs/utils"
)

func StartMaster() {
	fmt.Println("Master running")
	cm := ClientManager{
		register:   make(chan *Node),
		unregister: make(chan *Node),
	}
	/*******
    wm := WorkerManager{
		workers:    list.New(),
		pworker:    nil,
		register:   make(chan *Node),
		unregister: make(chan *Node),
	}
    **********/
	go cm.run()
	go cm.listen(MasterAddrToC)
	/************
    go wm.run()
	go wm.listen(MasterAddrToW)
    *************/
}

