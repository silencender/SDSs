package main

import (
	. "github.com/silencender/SDSs/utils"
	"github.com/silencender/SDSs/client"
    "log"
)

func main() {
    log.Println("client running")
    c := client.NewClient(1) 
    c.StartClient()
    WaitForINT(func() {})
}
