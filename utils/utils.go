package utils

import (
	"net"
)

const (
	MasterAddr = "localhost:12345"
)

type Node struct {
	socket net.Conn
	ok     bool
	info   net.Addr
	data   chan []byte
}

func (node *Node) close() {
	if node.ok {
		node.ok = false
		close(node.data)
		node.socket.Close()
	}
}
