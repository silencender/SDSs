package worker

type Worker struct {
	addr string
}

func NewWorker(addr string) *Worker {
	worker := Worker{addr: addr}
	return &worker
}

func (worker *Worker) StartWorker() {
}
