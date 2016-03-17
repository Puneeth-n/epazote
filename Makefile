.PHONY: all get test clean build cover compile goxc bintray

VERSION=1.5.0
GO ?= go
BIN_NAME=epazote
GO_XC = ${GOPATH}/bin/goxc -os="freebsd openbsd netbsd solaris dragonfly darwin linux"
GOXC_FILE = .goxc.local.json
GITHASH=`git rev-parse HEAD`

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
	${GO} build -ldflags "-X main.version=${VERSION} -X main.githash=${GITHASH}" -o ${BIN_NAME} cmd/epazote/main.go; \
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

define GOXCJSON
{
    "ConfigVersion": "0.9",
    "AppName": "epazote",
    "ArtifactsDest": "build",
    "PackageVersion": "${VERSION}",
    "TaskSettings": {
        "bintray": {
            "downloadspage": "bintray.md",
            "package": "epazote",
            "repository": "epazote",
            "subject": "nbari"
        }
    },
    "BuildSettings": {
        "LdFlags": "-X main.version=${VERSION} gg-X main.githash=${GITHASH}"
    }
}
endef

export GOXCJSON

goxc:
	@echo "$$GOXCJSON" > .goxc.json
	# ${GO_XC}

bintray:
	${GO_XC} bintray
