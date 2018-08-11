VERSION=`git describe --dirty`
LDFLAGS=-ldflags "-X main.version=${VERSION}"

build:
	go build ${LDFLAGS} -o ./bin/robin ./cmd/robin
	go build ${LDFLAGS} -o ./bin/robintest ./cmd/robintest
	go test -c -o ./bin/robin.cover -covermode=count -coverpkg=./,./cmd/robin,./cmd/robintest ./cmd/robin
	go test -c -o ./bin/robintest.cover -covermode=count -coverpkg=./,./cmd/robintest,./cmd/robin ./cmd/robintest

install:
	go build ${LDFLAGS} -o ${GOPATH:=~/go}/bin/robin ./cmd/robin
	go build ${LDFLAGS} -o ${GOPATH:=~/go}/bin/robintest ./cmd/robintest
	go test -c -o ${GOPATH:=~/go}/bin/robin.cover -covermode=count -coverpkg=./,./cmd/robin,./cmd/robintest ./cmd/robin
	go test -c -o ${GOPATH:=~/go}/bin/robintest.cover -covermode=count -coverpkg=./,./cmd/robintest,./cmd/robin ./cmd/robintest
