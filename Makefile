.PHONY: all get test clean build

GO ?= go
BIN_NAME=epazote

all: clean build

get:
	${GO} get

build: get
	${GO} build -o ${BIN_NAME} cmd/epazote/main.go

clean:
	@rm -f ${BIN_NAME}

test: get
	${GO} test -v
