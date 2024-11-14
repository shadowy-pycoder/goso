APP_NAME=goso
TEST_FILE=${APP_NAME}.test

all: build

.PHONY: build
build: 
	CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ./bin/${APP_NAME} ./cmd/${APP_NAME}/*.go


.PHONY: bench
bench: 
	go test -bench=. -benchmem -run=^$$ -benchtime 100x -cpuprofile='cpu.prof' -memprofile='mem.prof'

.PHONY: demo
demo:
	go test . -v -run=TestOutput -count=1 