# crawler

Script to crawl html and add href links, crawling and indexing 5k sites / second into [big-data graph DB](https://github.com/dgoldstein1/graphApi).

[![Maintainability](https://api.codeclimate.com/v1/badges/0918dd40ac9fd5d3e454/maintainability)](https://codeclimate.com/github/dgoldstein1/crawler/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/0918dd40ac9fd5d3e454/test_coverage)](https://codeclimate.com/github/dgoldstein1/crawler/test_coverage)
[![CircleCI](https://circleci.com/gh/dgoldstein1/crawler.svg?style=svg)](https://circleci.com/gh/dgoldstein1/crawler)

## Build it

#### Binary

```sh
dep ensure -v
go install
```

#### Docker
```sh
docker build . -t dgoldstein1/wikipedia-path
```

## Run it

```sh
dc up -d
```

or with dependencies running locally

```sh
# run crawl on wikipedia
export GRAPH_DB_ENDPOINT="http://localhost:5000" # endpoint of graph database
export TWO_WAY_KV_ENDPOINT="http://localhost:5001" # endpoint of k:v <-> v:k lookup metadata db
export STARTING_ENDPOINT="https://en.wikipedia.org/wiki/String_cheese"
export METRICS_PORT=8002 # port where prom metrics are served
export MAX_APPROX_NODES=1000 # approximate number of nodes to visit (+/- one order of magnitude), set to '-1' for unlimited crawl
./build/crawler --help
```


## Development

#### Local Development

- Install [inotifywait](https://linux.die.net/man/1/inotifywait)
```sh
./watch_dev_changes.sh
```

#### Testing

```sh
go test $(go list ./... | grep -v /vendor/)
```

## Authors

* **David Goldstein** - [DavidCharlesGoldstein.com](http://www.davidcharlesgoldstein.com/?github-wikipeida-path) - [Decipher Technology Studios](http://deciphernow.com/)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
