.DEFAULT_GOAL := build
.PHONY: lint tidy build serve

lint:
	golangci-lint run

tidy:
	go mod tidy

build:
	go build -ldflags "-w -s" -o legit

serve:
	./legit
