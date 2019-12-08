package main

import (
	. "github.com/silencender/SDSs/utils"
	"github.com/silencender/SDSs/client"
    "log"
    "os"
)

func main() {
    f,_ := os.OpenFile("client_result", os.O_RDWR | os.O_CREATE, 0666)
    log.SetOutput(f)
    log.Println("client running")
    c := client.NewClient(1,3)
    c.StartClient()
    WaitForINT(func() {})
}
