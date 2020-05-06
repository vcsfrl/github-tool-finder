package github_tool_finder

import (
	"bytes"
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

func (this *SearchReader) Handle() {
	request, _ := http.NewRequest("POST", "", bytes.NewBuffer([]byte(this.buildQl())))
	this.client.Do(request)
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
