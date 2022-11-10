go:
	go get -u google.golang.org/grpc
	go get -u golang.org/x/net
	go get -u google.golang.org/genproto/
	go get -u google.golang.org/protobuf/

install: 
	go get -u github.com/klauspost/reedsolomon
	go get -u github.com/klauspost/cpuid
	go get -u github.com/cbergoon/merkletree

build:
	export GOPATH=$PWD
	export GOBIN=$PWD/bin
	go install src/main/server.go
	go install src/main/client.go
	go install src/main/keygen.go 

all:
	go install build
	