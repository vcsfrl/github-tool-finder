package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	github_tool_finder "github.com/vcsfrl/github-tool-finder"
)

const repoSearchQuery = `{
  search(query: "language:go sort:stars-desc", type: REPOSITORY, first: 10) {
    repositoryCount
    edges {
      node {
        ... on Repository {
          description
          name
          nameWithOwner
          url
          owner {
            login
          }
          forkCount
          stargazers {
            totalCount
          }
          watchers {
            totalCount
          }
          homepageUrl
          licenseInfo {
            name
          }
          mentionableUsers {
            totalCount
          }
          mirrorUrl
          isMirror
          primaryLanguage {
            name
          }
          parent {
            name
          }
          createdAt
          updatedAt
        }
      }
    }
  }
}`

func main() {
	argLength := len(os.Args[1:])
	if argLength != 2 {
		printUsage()
	}

	client := github_tool_finder.NewAuthenticationClientV4(http.DefaultClient, "08ca17a189702d1a515f3218788944460ad0690c")
	request, _ := http.NewRequest("POST", "", bytes.NewBuffer([]byte(repoSearchQuery)))
	response, _ := client.Do(request)
	res, _ := ioutil.ReadAll(response.Body)
	log.Fatal(request)
	log.Fatal(string(res))

	//output := make(chan *github_tool_finder.Repository, 1024*1024)
	//nr, _ := strconv.Atoi(os.Args[2])
	//
	//reader := github_tool_finder.NewSearchReader(os.Args[1], nr, output, client)
	//go func() {
	//err := reader.Handle()
	//if nil != err {
	//	log.Fatal(err)
	//}
	//}()
	//writer := github_tool_finder.NewWriterHandler(output, os.Stdout)
	//writer.Handle()
}

func printUsage() {
	fmt.Print(`Usage:
  search [query] [nr]
`)
}
