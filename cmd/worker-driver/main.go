package main

import (
	. "github.com/silencender/SDSs/utils"
	"github.com/silencender/SDSs/worker"
)

func main() {
	w := worker.NewWorker()
	w.StartWorker()
	WaitForINT(func() {})
}
