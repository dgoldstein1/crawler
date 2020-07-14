FROM golang:1.14

# setup go
ENV GOBIN $GOPATH/bin
ENV PATH $GOBIN:/usr/local/go/bin:$PATH


COPY . /temp
RUN (cd /temp && go build -o build/crawler && cp build/crawler $GOBIN/crawler)
RUN rm -rf /temp


RUN crawler --version

# copy in assets
COPY synonyms/english.txt /english.txt
COPY ar_synonyms/arabic.txt /arabic.txt
COPY counties/counties.txt /counties.txt
ENV ENGLISH_WORD_LIST_PATH /english.txt
ENV ARABIC_WORD_LIST_PATH /arabic.txt
ENV COUNTIES_LIST /counties.txt


COPY LICENSE /LICENSE
COPY VERSION /VERSION

ENV COMMAND "--help"
CMD crawler $COMMAND
