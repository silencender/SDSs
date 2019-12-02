package client

import (
	"net"
)

type Client struct {
	socket net.Conn
}