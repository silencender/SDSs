package main

import (
	. "github.com/silencender/SDSs/utils"
	"github.com/silencender/SDSs/worker"
)

func main() {
    w := worker.NewWorker("localhost:18888")
	w.StartWorker()
	WaitForINT(func() {})
}
