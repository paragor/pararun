LINUX_BINARY_PATH=build/linux
PARARUN_MAIN=cmd/pararun/main.go

.PHONY: all
all: build run_remote

.PHONY: build
build:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o ${LINUX_BINARY_PATH} ${PARARUN_MAIN}

.PHONY: run_remote
run_remote:
	scp ${LINUX_BINARY_PATH} centos:/tmp/linux
	ssh centos sudo /tmp/linux


