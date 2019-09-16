package main

import (
	"github.com/dgoldstein1/crawler/crawler"
	db "github.com/dgoldstein1/crawler/db"
	wiki "github.com/dgoldstein1/crawler/wikipedia"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"strconv"
)

// checks environment for required env vars
var logFatalf = log.Fatalf
var logMsg = log.Infof

func parseEnv() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	requiredEnvs := []string{
		"GRAPH_DB_ENDPOINT",
		"MAX_APPROX_NODES",
		"TWO_WAY_KV_ENDPOINT",
		"METRICS_PORT",
	}
	for _, v := range requiredEnvs {
		if os.Getenv(v) == "" {
			logFatalf("'%s' was not set", v)
		} else {
			// print out config
			logMsg("%s=%s", v, os.Getenv(v))
		}
	}
	numberVars := []string{"MAX_APPROX_NODES", "PARALLELISM", "MS_DELAY"}
	for _, e := range numberVars {
		i, err := strconv.Atoi(os.Getenv(e))
		if err != nil {
			logFatalf("Could not parse %s for env variable %s. Reccieve: %v", e, os.Getenv(e), err.Error())
		}
		if i < 1 && i != -1 {
			logFatalf("%s must be greater than 1 but was '%i'", e, i)
		}

	}
}

// runs crawler with given functions
func runCrawler(
	isValidCrawlLink crawler.IsValidCrawlLinkFunction,
	addEdgeIfDoesNotExist crawler.AddEdgeFunction,
	getNewNode crawler.GetNewNodeFunction,
) {
	// assert environment
	parseEnv()
	// crawl with passed args
	crawler.ServeMetrics()
	crawler.Run(
		os.Getenv("STARTING_ENDPOINT"),
		isValidCrawlLink,
		db.ConnectToDB,
		addEdgeIfDoesNotExist,
		getNewNode,
	)
}

func main() {
	app := cli.NewApp()
	app.Name = "crawler"
	app.Usage = " acustomizable web crawler script for different websites"
	app.Description = "web crawl different URLs and add similar urls to a graph database"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:    "wikipedia",
			Aliases: []string{"w"},
			Usage:   "crawl on wikipedia articles",
			Action: func(c *cli.Context) error {
				runCrawler(
					wiki.IsValidCrawlLink,
					wiki.AddEdgesIfDoNotExist,
					wiki.GetRandomArticle,
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
