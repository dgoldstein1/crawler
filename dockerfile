FROM golang:1.9

# setup go
ENV GOBIN $GOPATH/bin
ENV PATH $GOBIN:/usr/local/go/bin:$PATH

COPY build $GOBIN
RUN crawler --version

# copy in assets
COPY synonyms/english.txt /english.txt
ENV ENGLISH_WORD_LIST_PATH /english.txt

ENV COMMAND "--help"
CMD crawler $COMMAND
