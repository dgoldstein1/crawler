# wikipedia-path

Finds path between two wikipedia articles.

## Required environment variables

```
- "NEO4J_ENDPOINT"
- "NEO4J_USERNAME"
- "NEO4J_PASS"
```

## Run it

```sh
dep ensure -v
fresh
```

or
```sh
docker build . -t dgoldstein1/wikipedia-path
docker run -p 8080:8080 dgoldstein1/wikipedia-path
```

## Test

```sh
go test $(go list ./... | grep -v /vendor/)
```

or

```sh
run_tests.sh
```

## API

`/metrics` -- shows prometheus metrics for the service

`/path?aricle1=Test_(wrestler)&article2=The_Un-Americans` -- generates a random number with a max

```json
{
	"edgeNodes" : {
		"article1" : {
			"valid" : true,
			"description" : "...",
		},
		"article2" : {
			"valid" : true,
			"description" : "..."
		}
	},
	"path" : {
		"exists" : true,
		"nodes" : [
			"Test_(wrestler)",
			"World_Tag_Team_Championship_(WWE)",
			"The_Un-Americans",
		],
	}

}
```


## Authors

* **David Goldstein** - [DavidCharlesGoldstein.com](http://www.davidcharlesgoldstein.com/?github-password-service) - [Decipher Technology Studios](http://deciphernow.com/)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
