.PHONY: all get test clean build cover

GO ?= go
BIN_NAME=epazote

all: clean build

get:
	${GO} get -u

build: get
	${GO} build -o ${BIN_NAME} cmd/epazote/main.go

clean:
	@rm -f ${BIN_NAME} *.out

test: get
	${GO} test -v

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
