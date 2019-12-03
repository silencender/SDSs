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
		message := make([]byte, 1024)
		length, err := client.Socket.Read(message)
		if err != nil {
			cm.unregister <- client
			break
		}
		if length > 0 {
			fmt.Printf("Received from client %s: %s", client.Info.String(), string(message))
			client.Data <- []byte("Master response: " + string(message))
		}
	}

}

func (cm *ClientManager) send(client *Node) {
	for {
		select {
		case message, ok := <-client.Data:
			if !ok {
				return
			}
			client.Socket.Write(message)
		}
	}
}

func (cm *ClientManager) listen() {
	listener, err := net.Listen("tcp", MasterAddrToC)
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		client := &Node{
			Socket: conn,
			Ok:     false,
			Info:   conn.RemoteAddr(),
			Data:   make(chan []byte),
		}
		cm.register <- client
		go cm.receive(client)
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
