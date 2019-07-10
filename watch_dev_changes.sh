#!/bin/bash
while true; do

inotifywait -e modify,create,delete -r ./ && \
	clear
	go build ./... \
		&& go test ./... -v
done
