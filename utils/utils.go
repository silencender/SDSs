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
	MaxQ          = 16
)

type Node struct {
	Socket     net.Conn
	Ok         bool
	Info       net.Addr
	ReqData    chan []byte
	ResData    chan []byte
	ListenAddr string
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
		node.Socket.Close()
	}
}

type Payload struct {
	data   []byte
	length byte
}

func NewPayload() *Payload {
	return &Payload{
		data:   []byte{},
		length: 0,
	}
}

func (p *Payload) Load(data []byte) {
	p.data = data
	p.length = byte(len(data))
}

func (p *Payload) Encode() []byte {
	return prependByte(p.data, p.length)
}

func (p *Payload) Decode() []byte {
	return p.data
}

type PayloadParser struct {
	payloads []*Payload
	data     []byte
	length   int
}

func NewPayloadParser() *PayloadParser {
	pp := &PayloadParser{
		payloads: make([]*Payload, MaxQ),
		data:     make([]byte, BufSize),
		length:   0,
	}
	for i := range pp.payloads {
		pp.payloads[i] = NewPayload()
	}

	return pp
}

func (pp *PayloadParser) Parse(data []byte) []*Payload {
	copy(pp.data[pp.length:], data)
	pp.length += len(data)
	num, idx := 0, 0
	for ; num < MaxQ && idx < pp.length; num++ {
		l := int(pp.data[idx])
		pp.payloads[num].Load(pp.data[idx+1 : idx+1+l])
		idx += l + 1
	}
	pp.length -= idx

	return pp.payloads[:num]
}

func prependByte(x []byte, y byte) []byte {
	x = append(x, 0)
	copy(x[1:], x)
	x[0] = y
	return x
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
