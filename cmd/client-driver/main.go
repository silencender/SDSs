package main

import (
	"log"

	"github.com/silencender/SDSs/client"
	. "github.com/silencender/SDSs/utils"
)

func main() {
	log.Println("client running")
	c := client.NewClient(1, 3)
	c.StartClient()
	WaitForINT(func() {})
}
