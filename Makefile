.PHONY: all get test clean build cover

GO ?= go
BIN_NAME=epazote

all: clean build

get:
	${GO} get

build: get
ifdef DEBUG
# make build DEBUG=true
	${GO} get -u github.com/mailgun/godebug
	${GOPATH}/bin/godebug build -o ${BIN_NAME}.debug cmd/epazote/main.go
else
	${GO} build -o ${BIN_NAME} cmd/epazote/main.go
endif

clean:
	@rm -f ${BIN_NAME} ${BIN_NAME}.debug *.out

test: get
	${GO} test -v

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
