package main

import (
	"log"
	"os"

	"github.com/silencender/SDSs/master"
	. "github.com/silencender/SDSs/utils"
)

func main() {
	f, _ := os.OpenFile("master_result", os.O_RDWR|os.O_CREATE, 0666)
	log.SetOutput(f)
	m := master.NewMaster("127.0.0.1:12345", "127.0.0.1:12346")
	m.StartMaster()
	WaitForINT(func() {})
}
