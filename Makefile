VERSION=1.0.0

default: build

build: test
	@echo Building server code
	@mkdir -p ./bin
	@go build ./socks5
	@go build ./proxy
	@go build ./handler
	@echo Building binary
	@mkdir -p ./bin
	@echo Building binary version $(VERSION)
	@go build -o ./bin/proxy -ldflags "-X main.version=$(VERSION)" ./cmd/proxy_main.go

test:
	@echo Executing unit tests
	@go test ./proxy

clean:
	@echo Cleaning up binaries
	@rm -rf bin