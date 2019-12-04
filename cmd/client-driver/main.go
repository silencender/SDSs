package main

import (
    "github.com/silencender/SDSs/client"
    "fmt"
)

func main() {
    fmt.Println("Hello,client")
    client := client.StartClient()
    if client.Master.Socket == nil {
        return
    }
    client.Run("i","1","2")
}
