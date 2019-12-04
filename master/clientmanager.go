package master

import (
	"fmt"
	"net"

	. "github.com/silencender/SDSs/utils"
)

type ClientManager struct {
	register   chan *Node
	unregister chan *Node
}

func (cm *ClientManager) receive(client *Node) {
	for {
		message := make([]byte, BufSize)
		length, err := client.Socket.Read(message)
		if err != nil {
			cm.unregister <- client
			break
		}
		if length > 0 {
			fmt.Printf("Received from client %s: %s", client.Info.String(), string(message))
			client.ReqData <- message
		}
	}

}

func (cm *ClientManager) handle(client *Node) {
	for {
		select {
		case req := <-client.ReqData:
			client.ResData <- []byte("Master response: " + string(req))
		}
	}
}

func (cm *ClientManager) send(client *Node) {
	for {
		select {
		case message, ok := <-client.ResData:
			if !ok {
				return
			}
			client.Socket.Write(message)
		}
	}
}

func (cm *ClientManager) listen(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		client := NewNode(conn)
		cm.register <- client
		go cm.receive(client)
		go cm.handle(client)
		go cm.send(client)
	}
}

func (cm *ClientManager) run() {
	for {
		select {
		case conn := <-cm.register:
			conn.Open()
			fmt.Printf("Client %s registered\n", conn.Info.String())
		case conn := <-cm.unregister:
			conn.Close()
			fmt.Printf("Client %s unregistered\n", conn.Info.String())
		}
	}
}
