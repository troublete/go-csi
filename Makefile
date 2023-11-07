build:
	mkdir -p ./dist/
	go build -o ./dist/build ./cmd/build/...
	go build -o ./dist/watch ./cmd/watch/...
	go build -o ./dist/serve ./cmd/serve/...

