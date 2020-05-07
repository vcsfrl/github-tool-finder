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
		return
	}

	token, ok := os.LookupEnv("GH_TOKEN")
	if !ok {
		fmt.Println("Please specify a github token (env variable: GH_TOKEN)")
		return
	}

	client := github_tool_finder.NewAuthenticationClientV4(http.DefaultClient, token)
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
