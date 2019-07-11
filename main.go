package wikipedia

import(
  "strings"
  "os"
  "log"
  "fmt"
  "io"
  "github.com/urfave/cli"
)

// checks environment for required env vars
var logFatalf = log.Fatalf
func parseEnv() {
  requiredEnvs := []string{
    "GRAPH_DB_ENDPOINT",
    "STARTING_ENDPOINT",
    "MAX_CRAWL_DEPTH",
  }
  for _, v := range requiredEnvs{
    if (os.Getenv(v) == "") {
      logFatalf("'%s' was not set", v)
    }
  }
}

func isValidCrawlLink(link string) bool {
  return strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":")
}

func main() {

    app := cli.NewApp()
    app.Name = "wikipedia_cralwer"
    app.Usage = "crawl wikipedia adding link articles into a graph database"
    app.Version = "0.1.0"
    // EXAMPLE: Append to an existing template
    cli.AppHelpTemplate = fmt.Sprintf(`%s

  WEBSITE: http://davidcharlesgoldstein.com

  SUPPORT: david0124816@gmail.com

  `, cli.AppHelpTemplate)

    cli.AppHelpTemplate = `NAME:
     {{.Name}} - {{.Usage}}
  USAGE:
     {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
     {{if len .Authors}}
  AUTHOR:
     David Goldstein
  REQUIRED ENV VARS:
    - "GRAPH_DB_ENDPOINT" - endpoint of the graphDB
    - "STARTING_ENDPOINT"- starting endpoint to crawl from
    - "MAX_CRAWL_DEPTH" - maximum depth to crawl
  COPYRIGHT:
     MIT
  VERSION:
     {{.Version}}
     {{end}}
  `

    // EXAMPLE: Replace the `HelpPrinter` func
    cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
      fmt.Println("Ha HA.  I pwnd the help!!1")
    }

    err := cli.NewApp().Run(os.Args)
    if err != nil {
      log.Fatal(err)
    }


}