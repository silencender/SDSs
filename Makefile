all:
	git add *
	git commit -m "auto-shell"
	git push
	cd ../github.com/silencender/SDSs
	git pull
	#protoc --go_out=.  protos/*.proto
	go build -o bin/client ./cmd/client-driver
	go build -o bin/worker ./cmd/worker-driver
	go build -o bin/master ./cmd/master-driver


clean:
	rm bin/*
