go:
	go get -u google.golang.org/grpc
	go get -u golang.org/x/net
	go get -u google.golang.org/genproto/

install:
	go get -u github.com/klauspost/reedsolomon
	go get -u github.com/klauspost/cpuid
	go get -u github.com/cbergoon/merkletree
	go get -u github.com/golang/protobuf

build:
	go install src/main/server.go
	go install src/main/client.go
	go install src/main/keygen.go

