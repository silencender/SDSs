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
	master.StartMaster()
	WaitForINT(func() {})
}
