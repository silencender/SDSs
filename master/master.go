package master

import (
	"fmt"

	. "github.com/silencender/SDSs/utils"
)

func StartMaster() {
	fmt.Println("Master running")
	cm := ClientManager{
		register:   make(chan *Node),
		unregister: make(chan *Node),
	}
	go cm.run()
	cm.listen()
}
