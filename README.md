# wikipedia-path

Script to crawl wikipedia articles and add them to them to a [big-data graph DB](https://github.com/dgoldstein1/graphApi)

## Run it

```sh
./crawler --endpoint "http://localhost:6080"
```

## Build it

#### Binary

```sh
go build
```

#### Docker
```sh
docker build . -t dgoldstein1/wikipedia-path
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
