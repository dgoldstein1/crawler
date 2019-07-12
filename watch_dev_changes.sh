#!/bin/bash
while true; do

inotifywait -e modify,create,delete -r ./ && \
	clear
	go build -o build/crawler \
		&& go test ./... -v
done
