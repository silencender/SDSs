package master

import (
	. "github.com/silencender/SDSs/utils"
)

type ClientManager struct {
	register   chan *Node
	unregister chan *Node
}

func (cm *ClientManager) receive(client *Node) {

}

func (cm *ClientManager) send(client *Node) {

}

func (cm *ClientManager) run() {

}
