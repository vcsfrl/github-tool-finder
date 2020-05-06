package github_tool_finder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type SearchReader struct {
	output chan *Repository
	client HTTPClient
	query  string
	number int
}

func (this *SearchReader) Close() error {
	close(this.output)

	return nil
}

func (this *SearchReader) Handle() error {
	defer this.Close()
	request, _ := http.NewRequest("POST", "", bytes.NewBuffer([]byte(this.buildQl())))
	response, _ := this.client.Do(request)
	decoder := json.NewDecoder(response.Body)
	result := &SearchResponse{}
	decoder.Decode(result)

	for _, edge := range result.Data.Search.Edges {
		this.output <- &edge.Node
	}

	return response.Body.Close()

}

func (this *SearchReader) buildQl() string {
	return fmt.Sprintf(repoSearchQuery, this.query, this.number)
}

func NewSearchReader(query string, number int, output chan *Repository, client HTTPClient) *SearchReader {
	return &SearchReader{query: query, number: number, output: output, client: client}
}

const repoSearchQuery = `{
  search(query: "%s", type: REPOSITORY, first: %d) {
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
