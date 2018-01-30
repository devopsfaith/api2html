.PHONY: all prepare deps test build server_build docker

GOLANG_VERSION=1.9.3-alpine3.7
DEP_VERSION=0.4.1

all: deps test build

prepare:
	@echo "Installing dep..."
	curl -L https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-${OS}-amd64 -o ${GOPATH}/bin/dep
	ls -l ${GOPATH}/bin/dep
	chmod a+x ${GOPATH}/bin/dep

deps:
	@echo "Setting up the vendors folder..."
	@dep ensure -v
	@echo ""
	@echo "Resolved dependencies:"
	@dep status
	@echo ""

test:
	go test -cover -v ./...

build:
	@echo "Building the binary..."
	@go build -a -o api2html
	@echo "You can now use ./api2html"

docker: server_build
	docker build -t devopsfaith/api2html .
	rm api2html

server_build: deps
	docker run --rm -it -e "GOPATH=/go" -v "${PWD}:/go/src/github.com/devopsfaith/api2html" -w /go/src/github.com/devopsfaith/api2html golang:${GOLANG_VERSION} go build -o api2html
