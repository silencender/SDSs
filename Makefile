all:
	go build -o bin/client ./cmd/client-driver
	go build -o bin/worker ./cmd/worker-driver
	go build -o bin/master ./cmd/master-driver


clean:
	rm bin/*
