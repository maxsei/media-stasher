all: build
fmt:
	go run golang.org/x/tools/cmd/goimports -w .
	go fmt
build: fmt
	go build -o ./bin/out *.go
	npm run build
run: build
	./bin/out
dev:
	npm run dev &
	find -name '*.go' | entr -sr 'go run *.go'
