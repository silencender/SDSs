package main

import (
	"log"
	"os"

	. "github.com/silencender/SDSs/utils"
	"github.com/silencender/SDSs/worker"
)

func main() {
	f, _ := os.OpenFile("worker_result", os.O_RDWR|os.O_CREATE, 0666)
	log.SetOutput(f)
	w1 := worker.NewWorker(18888)
	w2 := worker.NewWorker(18889)
	w3 := worker.NewWorker(18890)
	w1.StartWorker()
	w2.StartWorker()
	w3.StartWorker()
	WaitForINT(func() {})
}
