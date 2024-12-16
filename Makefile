export CGO_ENABLED=0

build:
	go build -o ./seekerbot .

run:
	./seekerbot server

.PHONY: build
