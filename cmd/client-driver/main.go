package main

import (
	"github.com/silencender/SDSs/client"
    "log"
)

func main() {
    log.Println("client running")
    client := client.StartClient()
    if client.Master.Socket == nil {
        return
    }
    client.Run("i","1","2")
}
