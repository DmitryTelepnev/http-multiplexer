GOOS ?= linux
GOARCH ?= amd64

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go get ./... && \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go mod tidy && \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go mod vendor && \
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
			-tags "netgo std static_all" \
    		-mod vendor \
    		-o ./bin/app cmd/main.go

vendor:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go get ./... && \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go mod tidy && \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go mod vendor

test:
	docker run -v ${PWD}/:/app -w /app golang:1.14-stretch sh -c \
	'go mod vendor && go test -mod vendor -bench=. -benchmem -v ./internal/...'