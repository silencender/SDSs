package main

import (
	"github.com/silencender/SDSs/master"
	. "SDSs/utils"
)

func main() {
	master.StartMaster()
	WaitForINT(func() {})
}

