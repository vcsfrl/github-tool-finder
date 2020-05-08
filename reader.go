package github_tool_finder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type SearchReader struct {
	client   HTTPClient
	query    string
	nrRepos  int
	pageSize int
	output   chan *Repository
}

func (this *SearchReader) Close() error {
	close(this.output)

	return nil
}

func (this *SearchReader) Handle() error {
	defer this.Close()
	result := this.readRepositories()

	return this.sendResult(result)
}

func (this *SearchReader) sendResult(result *SearchResponse) error {
	if "" != result.ErrorMessage {
		return errors.New(result.ErrorMessage)
	}

	for _, edge := range result.Data.Search.Edges {
		node := edge.Node
		this.output <- &node
	}

	return nil
}

func (this *SearchReader) readRepositories() *SearchResponse {
	result := &SearchResponse{}
	reader, err := this.repositoryReader()
	if nil != reader {
		defer reader.Close()
	}
	this.decodeRepositories(reader, result, err)

	return result
}

func (this *SearchReader) decodeRepositories(reader io.ReadCloser, result *SearchResponse, err error) {
	if nil != err {
		result.ErrorMessage = err.Error()
		return
	}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(result)
	if nil != err {
		result.ErrorMessage = err.Error()
	}
}

func (this *SearchReader) repositoryReader() (io.ReadCloser, error) {
	request, err := http.NewRequest("POST", "", strings.NewReader(this.buildQl()))
	if nil != err {
		return nil, err
	}
	response, err := this.client.Do(request)

	if nil != err {
		return nil, err
	}

	return response.Body, nil
}

func (this *SearchReader) buildQl() string {
	return fmt.Sprintf(repoSearchQuery, this.query, this.nrRepos)
}

func NewSearchReader(query string, number int, output chan *Repository, client HTTPClient) *SearchReader {
	return &SearchReader{query: query, nrRepos: number, output: output, client: client, pageSize: 100}
}

const repoSearchQuery = "{\"query\":\"query SearchRepositories {\\n  search(query: \\\"%s\\\", type: REPOSITORY, first:%d){\\n    repositoryCount\\n    edges {\\n      cursor \\n      node {\\n\\t\\t\\t\\t... on Repository {\\n          description\\n          name\\n          nameWithOwner\\n          url\\n          owner {\\n            login\\n          }\\n          forkCount\\n          stargazers {\\n            totalCount\\n          }\\n          watchers {\\n            totalCount\\n          }\\n          homepageUrl\\n          licenseInfo {\\n            name\\n          }\\n          mentionableUsers {\\n            totalCount\\n          }\\n          mirrorUrl\\n          isMirror\\n          primaryLanguage {\\n            name\\n          }\\n          parent {\\n            name\\n          }\\n          createdAt\\n          updatedAt\\n        }\\n      }\\n    }\\n  }\\n}\\n\",\"variables\":{}}"
