package main

import (
	"github.com/dgoldstein1/crawler/crawler"
	wiki "github.com/dgoldstein1/crawler/wikipedia"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
)

// checks environment for required env vars
var logFatalf = log.Fatalf

func parseEnv() {
	requiredEnvs := []string{
		"GRAPH_DB_ENDPOINT",
		"STARTING_ENDPOINT",
		"MAX_CRAWL_DEPTH",
	}
	for _, v := range requiredEnvs {
		if os.Getenv(v) == "" {
			logFatalf("'%s' was not set", v)
		}
	}
	i, err := strconv.Atoi(os.Getenv("MAX_CRAWL_DEPTH"))
	if err != nil {
		logFatalf(err.Error())
	}
	if i < 1 {
		logFatalf("MAX_CRAWL_DEPTH must be greater than 1 but was '%i'", i)
	}
}

// runs crawler with given functions
func runCrawler(
	isValidCrawlLink crawler.IsValidCrawlLinkFunction,
	connectToDB crawler.ConnectToDBFunction,
	addEdgeIfDoesNotExist crawler.AddEdgeFunction,
) {
	// assert environment
	parseEnv()
	// crawl with passed args
	MAX_CRAWL_DEPTH, _ := strconv.Atoi(os.Getenv("MAX_CRAWL_DEPTH"))
	crawler.Crawl(
		os.Getenv("STARTING_ENDPOINT"),
		int32(MAX_CRAWL_DEPTH),
		isValidCrawlLink,
		connectToDB,
		addEdgeIfDoesNotExist,
	)
}

func main() {
	app := cli.NewApp()
	app.Name = "crawler"
	app.Usage = " acustomizable web crawler script for different websites"
	app.Description = "web crawl different URLs and add similar urls to a graph database"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "wikipedia",
			Aliases: []string{"w"},
			Usage:   "crawl on wikipedia articles",
			Action: func(c *cli.Context) error {
				os.Setenv("WIKI_API_ENDPOINT", "https://en.wikipedia.org/w/api.php")
				runCrawler(
					wiki.IsValidCrawlLink,
					wiki.ConnectToDB,
					wiki.AddEdgesIfDoNotExist,
				)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
