package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	http2 "github.com/vcsfrl/github-tool-finder/http"

	"github.com/vcsfrl/github-tool-finder/search"
)

func main() {
	query, total, token := getArguments()

	transport := make(chan *search.Repository, 1024*1024)
	client := http2.NewAuthenticationClientV4(http.DefaultClient, token)
	reader := search.NewRepositoryReader(query, total, transport, client)
	writer := search.NewCsvWriter(transport, os.Stdout)

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
		fmt.Fprintln(os.Stderr, usage())
		os.Exit(1)
	}

	token, ok := os.LookupEnv("GH_TOKEN")
	if !ok {
		fmt.Fprintln(os.Stderr, "Please specify a github token (environment variable: GH_TOKEN).")
		os.Exit(1)
	}

	total, err := strconv.Atoi(os.Args[2])
	if nil != err {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return os.Args[1], total, token
}

func usage() string {
	return fmt.Sprintf(`
Usage:
 search [query] [total]

`)
}
