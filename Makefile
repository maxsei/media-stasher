all: build
build:
	echo "Building binary"
	npm run build
	go build -o ./bin/out *.go
run: build
	echo "Running binary"
	./bin/out
