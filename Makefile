PROGRAM_NAME=media-stasher

all: build
fmt:
	go generate ./...
	go run golang.org/x/tools/cmd/goimports -w .
	go fmt
build: fmt
	go build -o ./bin/$(PROGRAM_NAME) $(filter-out ./tools.go, $(wildcard ./*.go))
	npm run build
run: build
	./bin/$(PROGRAM_NAME)
debug: build
	go run github.com/go-delve/delve/cmd/dlv exec ./bin/$(PROGRAM_NAME)
dev: fmt
	npm run dev &
	find -name '*.go' | entr -sr ' go run *.go'
clean:
	go clean -testcache -cache
	rm -f $(filter-out ./bin/.gitkeep, $(wildcard ./bin/*))
	rm -f $(wildcard ./public/build/*)
