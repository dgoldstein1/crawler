FROM golang:latest

ENV PRJ_DIR $GOPATH/src/github.com/dgoldstein1/cmd/wikipedia-path
# create project dir
RUN mkdir -p $PRJ_DIR
# add src, service communication ,and docs
COPY . $PRJ_DIR
RUN mkdir -p mkdir -p /opt/services/wikipedia-path
COPY ./Gopkg.toml $PRJ_DIR
COPY ./Gopkg.lock $PRJ_DIR

# setup go
ENV GOBIN $GOPATH/bin
ENV PATH $GOBIN:/usr/local/go/bin:$PATH

# install utils
RUN go get github.com/golang/dep/cmd/dep

# copy over source code
WORKDIR $PRJ_DIR

# install dependencies
RUN dep ensure -v

# configure entrypoint
RUN go build

ENTRYPOINT ["./wikipedia-path"]

# expose service ports
EXPOSE 10000
EXPOSE 10001
EXPOSE 8080