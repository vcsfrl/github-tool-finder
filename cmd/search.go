package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	github_tool_finder "github.com/vcsfrl/github-tool-finder"
)

func main() {
	argLength := len(os.Args[1:])
	if argLength != 2 {
		printUsage()
	}

	client := github_tool_finder.NewAuthenticationClient(http.DefaultClient, "https", "api.github.com/graphql", "08ca17a189702d1a515f3218788944460ad0690c")
	output := make(chan *github_tool_finder.Repository, 1024*1024)
	nr, _ := strconv.Atoi(os.Args[2])
	reader := github_tool_finder.NewSearchReader(os.Args[1], nr, output, client)
	go func() {
		err := reader.Handle()
		if nil != err {
			log.Fatal(err)
		}
	}()
	writer := github_tool_finder.NewWriterHandler(output, os.Stdout)
	writer.Handle()
}

func printUsage() {
	fmt.Print(`Usage:
  search [query] [nr]
`)
}
