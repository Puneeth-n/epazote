.PHONY: all get test clean build cover compile goxc bintray

VERSION=1.4.0
GO ?= go
BIN_NAME=epazote
GO_XC = ${GOPATH}/bin/goxc -os="freebsd openbsd netbsd solaris dragonfly darwin linux" -build-ldflags="-X main.version=${VERSION}"
GOXC_FILE = .goxc.local.json

all: clean build

get:
	${GO} get

build: get
# make build DEBUG=true
	@if test -n "${DEBUG}"; then \
	${GO} get -u github.com/mailgun/godebug; \
	${GOPATH}/bin/godebug build -instrument="github.com/nbari/epazote/..." -o ${BIN_NAME}.debug cmd/epazote/main.go; \
	else \
	${GO} get -u gopkg.in/yaml.v2; \
	${GO} build -ldflags "-X main.version=${VERSION} -X main.githash=`git rev-parse HEAD`" -o ${BIN_NAME} cmd/epazote/main.go; \
	fi;

clean:
	@rm -rf ${BIN_NAME} ${BIN_NAME}.debug *.out build debian

test: get
	${GO} test -v

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out

compile: goxc

goxc:
	$(shell sed -i '' -e 's/"PackageVersion.*/"PackageVersion": "${VERSION}",/g' .goxc.json)
	$(shell echo '{\n "ConfigVersion": "0.9",' > $(GOXC_FILE))
	$(shell echo ' "TaskSettings": {' >> $(GOXC_FILE))
	$(shell echo '  "bintray": {\n   "apikey": "$(BINTRAY_APIKEY)"' >> $(GOXC_FILE))
	$(shell echo '  }\n } \n}' >> $(GOXC_FILE))
	${GO_XC}

bintray:
	${GO_XC} bintray
