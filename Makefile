APP_NAME=goso
TEST_FILE=${APP_NAME}.test

all: build

.PHONY: build
build: 
	CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ./bin/${APP_NAME} ./cmd/${APP_NAME}/*.go


${TEST_FILE}:
	go test -c -race

.PHONY: bench
bench: ${TEST_FILE}
	./$< -test.bench=. -test.benchmem -test.run=^$$ -test.benchtime 1000x \
	-test.cpuprofile='cpu.prof' -test.memprofile='mem.prof'