all:
	protoc --go_out=.  protos/*.proto
	go build -o bin/client ./cmd/client_driver
	go build -o bin/worker ./cmd/worker_driver
	go build -o bin/master ./cmd/master_driver


clean:
	rm bin/*