package worker

import (
	"net"
)

type Worker struct {
	socket net.Conn
}