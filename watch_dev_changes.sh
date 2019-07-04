#!/bin/bash
while true; do

run_tests() {
	go test $(go list ./... | grep -v /vendor/)
}

inotifywait -e modify,create,delete -r ./ && \
	clear
	go build \
		&& run_tests
done
