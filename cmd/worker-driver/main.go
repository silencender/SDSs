package main

import (
	. "github.com/silencender/SDSs/utils"
	"github.com/silencender/SDSs/worker"
)

func main() {
	port := 12330
	//addr := "localhost:" + string(port)
	w := worker.NewWorker(string(port))
	w.StartWorker()
	WaitForINT(func() {})
}
