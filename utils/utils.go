package utils

import (
	"net"
)

const (
	MasterAddr = "localhost:12345"
)

type Node struct {
	Socket net.Conn
	Ok     bool
	Info   net.Addr
	Data   chan []byte
}

func (node *Node) close() {
	if node.Ok {
		node.Ok = false
		close(node.Data)
	    node.Socket.Close()
	}
}
