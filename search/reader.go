package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var ErrRead = errors.New("read error")

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

func (sr *SearchReader) Close() error {
	close(sr.output)

	return nil
}

func (sr *SearchReader) Handle() error {
	defer sr.Close()
	sr.adjustPageSize()

	return sr.paginatedRead()
}

func (sr *SearchReader) adjustPageSize() {
	if sr.pageSize > sr.total {
		sr.pageSize = sr.total
	}
}

func (sr *SearchReader) paginatedRead() error {
	var (
		result *SearchResponse
		cursor string
	)

	for i := 0; i < sr.total; i += sr.pageSize {
		result = sr.readRepositories(sr.calculateLimit(i), sr.findCursor(result, cursor))
		if err := sr.sendResult(result); nil != err {
			return err
		}
	}

	return nil
}

func (sr *SearchReader) calculateLimit(readIndex int) int {
	limit := sr.pageSize
	if readIndex+sr.pageSize > sr.total {
		limit = sr.total - readIndex
	}

	return limit
}

func (sr *SearchReader) findCursor(previousResult *SearchResponse, cursor string) string {
	if nil != previousResult {
		cursor = previousResult.Data.Search.Edges[len(previousResult.Data.Search.Edges)-1].Cursor
	}

	return cursor
}

func (sr *SearchReader) readRepositories(limit int, cursor string) *SearchResponse {
	result := &SearchResponse{}
	reader, err := sr.repositoryQueryReader(sr.buildQl(limit, cursor))

	if nil != reader {
		defer reader.Close()
	}

	sr.decodeRepositories(reader, result, err)

	return result
}

func (sr *SearchReader) sendResult(result *SearchResponse) error {
	if done, err := sr.getErrors(result); done {
		return err
	}

	for _, edge := range result.Data.Search.Edges {
		node := edge.Node
		sr.output <- &node
	}

	return nil
}

func (sr *SearchReader) decodeRepositories(reader io.Reader, result *SearchResponse, err error) {
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

func (sr *SearchReader) repositoryQueryReader(query string) (io.ReadCloser, error) {
	request, _ := http.NewRequest("POST", "", strings.NewReader(query))
	response, err := sr.client.Do(request)

	if nil != err {
		return nil, err
	}

	return response.Body, nil
}

func (sr *SearchReader) buildQl(limit int, cursor string) string {
	if cursor != "" {
		cursor = fmt.Sprintf(", after: \\\"%s\\\"", cursor)
	}

	return fmt.Sprintf(repoSearchQuery, sr.query, limit, cursor)
}

func (sr *SearchReader) getErrors(result *SearchResponse) (bool, error) {
	if result.Message != "" {
		return true, fmt.Errorf("%s: %w", result.Message, ErrRead)
	}

	if 0 < len(result.Errors) {
		err := sr.wrapErrors(result)
		return true, err
	}

	return false, nil
}

func (sr *SearchReader) wrapErrors(result *SearchResponse) error {
	err := fmt.Errorf("api error: %w", ErrRead)

	for _, resultErr := range result.Errors {
		err = fmt.Errorf("%s - %s: %w", resultErr.Type, resultErr.Message, err)
	}

	return err
}

func NewSearchReader(query string, total int, output chan *Repository, client HTTPClient) *SearchReader {
	return &SearchReader{query: query, total: total, output: output, client: client, pageSize: 100}
}

const repoSearchQuery = "{\"query\":\"query SearchRepositories {\\n" +
	"  search(query: \\\"%s\\\", type: REPOSITORY, first:%d%s){\\n" +
	"    repositoryCount\\n" +
	"    edges {\\n" +
	"      cursor \\n" +
	"      node {\\n" +
	"\\t\\t\\t\\t... on Repository {\\n" +
	"          description\\n" +
	"          name\\n" +
	"          nameWithOwner\\n" +
	"          url\\n" +
	"          owner {\\n" +
	"            login\\n" +
	"          }\\n" +
	"          forkCount\\n" +
	"          stargazers {\\n" +
	"            totalCount\\n" +
	"          }\\n" +
	"          watchers {\\n" +
	"            totalCount\\n" +
	"          }\\n" +
	"          homepageUrl\\n" +
	"          licenseInfo {\\n" +
	"            name\\n" +
	"          }\\n" +
	"          mentionableUsers {\\n" +
	"            totalCount\\n" +
	"          }\\n" +
	"          mirrorUrl\\n" +
	"          isMirror\\n" +
	"          primaryLanguage {\\n" +
	"            name\\n" +
	"          }\\n" +
	"          parent {\\n" +
	"            name\\n" +
	"          }\\n" +
	"          createdAt\\n" +
	"          updatedAt\\n" +
	"        }\\n" +
	"      }\\n" +
	"    }\\n" +
	"  }\\n" +
	"}\\n" +
	"\",\"variables\":{}}"
