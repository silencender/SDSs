package main

import (
	"log"
	"os"

	"github.com/silencender/SDSs/client"
	. "github.com/silencender/SDSs/utils"
)

func main() {
	f, _ := os.OpenFile("client_result", os.O_RDWR|os.O_CREATE, 0666)
	log.SetOutput(f)
	log.Println("client running")
	c1 := client.NewClient(1, 1000000)
	c2 := client.NewClient(2, 1000000)
	c3 := client.NewClient(3, 1000000)
	c1.StartClient()
	c2.StartClient()
	c3.StartClient()
	WaitForINT(func() {})
}
