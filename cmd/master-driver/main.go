package main

import (
	"github.com/silencender/SDSs/master"
	. "github.com/silencender/SDSs/utils"
)

func main() {
	master.StartMaster()
	WaitForINT(func() {})
}
