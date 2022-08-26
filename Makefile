.PHONY: all examples
all: generate build test examples local-install

# Only needed to be done once to prepare for tests
generate:
	go generate ./... || true

build:
	go build -o bin/protoc-gen-elmer cmd/protoc-gen-elmer/main.go
	go build -o bin/protoc-gen-elmer-fuzzer cmd/protoc-gen-elmer-fuzzer/main.go
	go build -o bin/protoc-gen-elmer-twirp cmd/protoc-gen-elmer-twirp/main.go

test:
	go test ./...

examples:
	make -C examples

local-install:
	cp bin/protoc-gen-elmer* ~/bin
