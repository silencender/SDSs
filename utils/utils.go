package utils

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
)

const (
	MasterAddrToC = "localhost:12345"
	MasterAddrToW = "localhost:12346"
	BufSize       = 1024
)

type Node struct {
	Socket  net.Conn
	Ok      bool
	Info    net.Addr
	ReqData chan []byte
	ResData chan []byte
}

func NewNode(conn net.Conn) *Node {
	return &Node{
		Socket:  conn,
		Ok:      false,
		Info:    conn.RemoteAddr(),
		ReqData: make(chan []byte),
		ResData: make(chan []byte),
	}
}

func (node *Node) Open() {
	if !node.Ok {
		node.Ok = true
	}
}

func (node *Node) Close() {
	if node.Ok {
		node.Ok = false
		close(node.ReqData)
		close(node.ResData)
		node.Socket.Close()
	}
}

func WaitForINT(callback func()) {
	signalChan := make(chan os.Signal, 1)
	block := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("Keyboad interrupt. Doing cleaning jobs...")
		callback()
		close(block)
	}()
	<-block
}

func PrintIfErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func RemoveListItem(l *list.List, item interface{}) {
	for p := l.Front(); ; {
		if p.Value == item {
			l.Remove(p)
			break
		}
		if p == l.Back() {
			break
		}
		p = p.Next()
	}
}
