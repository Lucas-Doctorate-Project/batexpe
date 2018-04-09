build:
	go build -o ./bin/robin ./cmd/robin
	go build -o ./bin/robintest ./cmd/robintest
	go test -c -o ./bin/robintest.cover -covermode=count -coverpkg=./,./cmd/robintest,./cmd/robin ./cmd/robintest
	go test -c -o ./bin/robin.cover -covermode=count -coverpkg=./,./cmd/robin,./cmd/robintest ./cmd/robin

install:
	go build -o ${GOPATH:=~/go}/bin/robin ./cmd/robin
	go build -o ${GOPATH:=~/go}/bin/robintest ./cmd/robintest	
	go test -c -o ${GOPATH:=~/go}/bin/robintest.cover -covermode=count -coverpkg=./,./cmd/robintest,./cmd/robin ./cmd/robintest
	go test -c -o ${GOPATH:=~/go}/bin/robin.cover -covermode=count -coverpkg=./,./cmd/robin,./cmd/robintest ./cmd/robin
