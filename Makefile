build:
	GOOS=linux GOARCH=amd64 go get ./... && \
		GOOS=linux GOARCH=amd64 go mod tidy && \
		GOOS=linux GOARCH=amd64 go mod vendor && \
	GOOS=linux GOARCH=amd64 go build \
			-tags "netgo std static_all" \
    		-mod vendor \
    		-o ./bin/app cmd/main.go

vendor:
	GOOS=linux GOARCH=amd64 go get ./... && \
		GOOS=linux GOARCH=amd64 go mod tidy && \
		GOOS=linux GOARCH=amd64 go mod vendor

test:
	docker run -v ${PWD}/:/app -w /app golang:1.14-stretch sh -c \
	'go mod vendor && go test -mod vendor -bench=. -benchmem -v ./internal/...'