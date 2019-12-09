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
	BufSize = 1024
	MaxQ    = 256
)

const (
	MaxInt32 = int32(^uint32(0) >> 1)
)

type Node struct {
	Socket     net.Conn
	Ok         bool
	Info       net.Addr
	Window     chan byte
	ReqData    chan []byte
	ResData    chan []byte
	ListenAddr string
}

func NewNode(conn net.Conn) *Node {
	node := &Node{
		Socket:  conn,
		Ok:      false,
		Info:    conn.RemoteAddr(),
		Window:  make(chan byte, MaxQ),
		ReqData: make(chan []byte),
		ResData: make(chan []byte),
	}
	for i := 0; i < MaxQ; i++ {
		node.Window <- 0
	}

	return node
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

func (node *Node) Acquire() {
	<-node.Window
}

func (node *Node) Release() {
	node.Window <- 0
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
	data := make([]byte, len(p.data))
	copy(data, p.data)
	return data
}

type PayloadParser struct {
	payloads []*Payload
	data     []byte
	idx      int
	length   int
}

func NewPayloadParser() *PayloadParser {
	pp := &PayloadParser{
		payloads: make([]*Payload, MaxQ),
		data:     make([]byte, BufSize*2),
		idx:      0,
		length:   0,
	}
	for i := range pp.payloads {
		pp.payloads[i] = NewPayload()
	}

	return pp
}

func (pp *PayloadParser) Parse(data []byte) []*Payload {
	//log.Println("data length,idx, length: ", len(data), pp.idx, pp.length)
	copy(pp.data, pp.data[pp.idx:pp.idx+pp.length])
	copy(pp.data[pp.length:], data)
	pp.length += len(data)
	pp.idx = 0
	num, idx, l := 0, 0, 0
	for ; num < MaxQ; num++ {
		l = int(pp.data[idx])
		idx += l + 1
		if idx > pp.length {
			break
		}
		pp.payloads[num].Load(pp.data[idx-l : idx])
		pp.idx = idx
	}
	pp.length -= pp.idx

	return pp.payloads[:num]
}

type SeqGen struct {
	seq int32
}

func NewSeqGen() *SeqGen {
	return &SeqGen{-1}
}

func (s *SeqGen) GetSeq() int32 {
	s.seq += 1
	s.seq %= MaxInt32

	return s.seq
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
