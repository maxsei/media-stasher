all: build
fmt:
	go generate ./...
	go run golang.org/x/tools/cmd/goimports -w .
	go fmt
build: fmt
	go build -o ./bin/out *.go
	npm run build
run: build
	./bin/out
debug: build
	go run github.com/go-delve/delve/cmd/dlv exec ./bin/out
dev: fmt
	npm run dev &
	find -name '*.go' | entr -sr ' go run *.go'
