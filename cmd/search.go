package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	github_tool_finder "github.com/vcsfrl/github-tool-finder"
)

func main() {
	query, total, token := getArguments()

	transport := make(chan *github_tool_finder.Repository, 1024*1024)
	client := github_tool_finder.NewAuthenticationClientV4(http.DefaultClient, token)
	reader := github_tool_finder.NewSearchReader(query, total, transport, client)
	writer := github_tool_finder.NewWriterHandler(transport, os.Stdout)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := reader.Handle()
		if nil != err {
			log.Fatal(err)
		}
		wg.Done()
	}()

	writer.Handle()
	wg.Wait()
}

func getArguments() (string, int, string) {
	argLength := len(os.Args[1:])
	if argLength != 2 {
		log.Fatal(usage())
	}

	token, ok := os.LookupEnv("GH_TOKEN")
	if !ok {
		log.Fatal("Please specify a github token (environment variable: GH_TOKEN)")
	}

	total, err := strconv.Atoi(os.Args[2])
	if nil != err {
		log.Fatal(err.Error())
	}

	return os.Args[1], total, token
}

func usage() string {
	return fmt.Sprintf(`Usage:
 search [query] [total]
`)
}
