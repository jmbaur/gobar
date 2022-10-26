build:
	go build -o $out/gobar ./cmd/gobar

check:
	staticcheck ./...
	go test ./...

run:
	go run ./cmd/gobar
