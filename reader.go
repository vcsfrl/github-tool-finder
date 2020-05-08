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
	query    string
	total    int
	pageSize int
	client   HTTPClient
	output   chan *Repository
}

func (this *SearchReader) Close() error {
	close(this.output)

	return nil
}

func (this *SearchReader) Handle() error {
	defer this.Close()
	this.adjustPageSize()

	return this.paginatedRead()
}

func (this *SearchReader) paginatedRead() error {
	var result *SearchResponse
	var cursor string
	for i := 0; i < this.total; i = i + this.pageSize {
		result = this.readRepositories(this.calculateLimit(i), this.findCursor(result, cursor))
		err := this.sendResult(result)
		if nil != err {
			return err
		}
	}
	return nil
}

func (this *SearchReader) findCursor(previousResult *SearchResponse, cursor string) string {
	if nil != previousResult {
		cursor = previousResult.Data.Search.Edges[len(previousResult.Data.Search.Edges)-1].Cursor
	}
	return cursor
}

func (this *SearchReader) adjustPageSize() {
	if this.pageSize > this.total {
		this.pageSize = this.total
	}
}

func (this *SearchReader) calculateLimit(readIndex int) int {
	limit := this.pageSize
	if readIndex+this.pageSize > this.total {
		limit = this.total - readIndex
	}
	return limit
}

func (this *SearchReader) sendResult(result *SearchResponse) error {
	if err, done := this.getErrors(result); done {
		return err
	}

	for _, edge := range result.Data.Search.Edges {
		node := edge.Node
		this.output <- &node
	}

	return nil
}

func (this *SearchReader) getErrors(result *SearchResponse) (error, bool) {
	if "" != result.Message {
		return errors.New(result.Message), true
	}

	if 0 < len(result.Errors) {
		err := this.wrapErrors(result)
		return err, true
	}
	return nil, false
}

func (this *SearchReader) wrapErrors(result *SearchResponse) error {
	err := errors.New("api error")
	for _, resultErr := range result.Errors {
		err = fmt.Errorf("%v, %w", err, errors.New(fmt.Sprintf("%s: %s", resultErr.Type, resultErr.Message)))
	}
	return err
}

func (this *SearchReader) readRepositories(limit int, cursor string) *SearchResponse {
	result := &SearchResponse{}
	reader, err := this.repositoryReader(limit, cursor)
	if nil != reader {
		defer reader.Close()
	}
	this.decodeRepositories(reader, result, err)

	return result
}

func (this *SearchReader) decodeRepositories(reader io.ReadCloser, result *SearchResponse, err error) {
	if nil != err {
		result.Message = err.Error()
		return
	}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(result)
	if nil != err {
		result.Message = err.Error()
	}
}

func (this *SearchReader) repositoryReader(limit int, cursor string) (io.ReadCloser, error) {
	request, err := http.NewRequest("POST", "", strings.NewReader(this.buildQl(limit, cursor)))
	if nil != err {
		return nil, err
	}
	response, err := this.client.Do(request)

	if nil != err {
		return nil, err
	}

	return response.Body, nil
}

func (this *SearchReader) buildQl(limit int, cursor string) string {
	if "" != cursor {
		cursor = fmt.Sprintf(", after: \\\"%s\\\"", cursor)
	}
	return fmt.Sprintf(repoSearchQuery, this.query, limit, cursor)
}

func NewSearchReader(query string, number int, output chan *Repository, client HTTPClient) *SearchReader {
	return &SearchReader{query: query, total: number, output: output, client: client, pageSize: 100}
}

const repoSearchQuery = "{\"query\":\"query SearchRepositories {\\n  search(query: \\\"%s\\\", type: REPOSITORY, first:%d%s){\\n    repositoryCount\\n    edges {\\n      cursor \\n      node {\\n\\t\\t\\t\\t... on Repository {\\n          description\\n          name\\n          nameWithOwner\\n          url\\n          owner {\\n            login\\n          }\\n          forkCount\\n          stargazers {\\n            totalCount\\n          }\\n          watchers {\\n            totalCount\\n          }\\n          homepageUrl\\n          licenseInfo {\\n            name\\n          }\\n          mentionableUsers {\\n            totalCount\\n          }\\n          mirrorUrl\\n          isMirror\\n          primaryLanguage {\\n            name\\n          }\\n          parent {\\n            name\\n          }\\n          createdAt\\n          updatedAt\\n        }\\n      }\\n    }\\n  }\\n}\\n\",\"variables\":{}}"
