FROM golang:1.9

RUN go get github.com/golang/dep/cmd/dep
RUN dep ensure -v
RUN go build -v ./... -o build/crawler


# build executable
RUN go install  ./...
RUN crawler --version


CMD ["crawler","$COMMAND"]
