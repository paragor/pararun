LINUX_BINARY_PATH=build/linux
PARARUN_MAIN=cmd/pararun/main.go

REMOTE_SSH=centos
REMOTE_DIR_BIN=/tmp/linux

.PHONY: all
all: centos build upload run_remote

.PHONY: build
build:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o ${LINUX_BINARY_PATH} ${PARARUN_MAIN}

.PHONY: upload
upload:
	scp ${LINUX_BINARY_PATH} ${REMOTE_SSH}:${REMOTE_DIR_BIN}
	ssh ${REMOTE_SSH} sudo chmod +x ${REMOTE_DIR_BIN}

.PHONY: run_remote
run_remote:
	ssh -t ${REMOTE_SSH} sudo ${REMOTE_DIR_BIN}

.PHONY: cetnos
centos:
	cd runtime/centos && vagrant up

