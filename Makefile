APPNAME = grpc-world
BIN = $(GOPATH)/bin
GOCMD = /usr/local/go/bin/go
GOBUILD = $(GOCMD) build
GOINSTALL = $(GOCMD) install
GORUN = $(GOCMD) run
BINARY_UNIX = $(BIN)/$(APPNAME)
PID = .pid
HUB_ADDR = hub.kpaas.nsini.com
TAG = v0.0.01-test
NAMESPACE = app
PWD = $(shell pwd)

start-server:
	$(BIN)/$(APPNAME) -http-addr :8080 -grpc-addr :8081 & echo $$! > $(PID)

restart:
	@echo restart the app...
	@kill `cat $(PID)` || true
	$(BIN)/$(APPNAME) -http-addr :8080 -grpc-addr :8081 & echo $$! > $(PID)

install:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOINSTALL) -v

stop:
	@kill `cat $(PID)` || true

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

docker-build:
	docker build --rm -t $(APPNAME):$(TAG) .

run:
	GO111MODULE=on $(GORUN) ./cmd/main.go -http-addr :8080 -grpc-addr :8081