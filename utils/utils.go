package utils

import (
	"net"
)

const (
	MasterAddrToC = "localhost:12345"
	MasterAddrToW = "localhost:12346"
	BufSize       = 1024
)

type Node struct {
	Socket net.Conn
	Ok     bool
	Info   net.Addr
	Data   chan []byte
}

func (node *Node) Open() {
	if !node.Ok {
		node.Ok = true
	}
}

func (node *Node) Close() {
	if node.Ok {
		node.Ok = false
		close(node.Data)
		node.Socket.Close()
	}
}
