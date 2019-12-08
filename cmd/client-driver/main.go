package main

import (
	"github.com/silencender/SDSs/client"
    "log"
)

func main() {
    log.Println("client running")
    client.NewClient(1) 
}
