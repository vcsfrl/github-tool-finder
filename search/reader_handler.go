package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	finderhttp "github.com/vcsfrl/github-tool-finder/http"
)

var ErrRead = errors.New("read error")

type RepositoryReader struct {
	query    string
	total    int
	pageSize int
	client   finderhttp.Client
	output   chan *Repository
}

func (sr *RepositoryReader) Close() error {
	close(sr.output)

	return nil
}

func (sr *RepositoryReader) Handle() error {
	defer sr.Close()
	sr.adjustPageSize()

	return sr.paginatedRead()
}

func (sr *RepositoryReader) adjustPageSize() {
	if sr.pageSize > sr.total {
		sr.pageSize = sr.total
	}
}

func (sr *RepositoryReader) paginatedRead() error {
	var (
		result *Response
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

func (sr *RepositoryReader) calculateLimit(readIndex int) int {
	limit := sr.pageSize
	if readIndex+sr.pageSize > sr.total {
		limit = sr.total - readIndex
	}

	return limit
}

func (sr *RepositoryReader) findCursor(previousResult *Response, cursor string) string {
	if nil != previousResult {
		cursor = previousResult.Data.Search.Edges[len(previousResult.Data.Search.Edges)-1].Cursor
	}

	return cursor
}

func (sr *RepositoryReader) readRepositories(limit int, cursor string) *Response {
	result := &Response{}
	reader, err := sr.repositoryQueryReader(sr.buildQl(limit, cursor))

	if nil != reader {
		defer reader.Close()
	}

	sr.decodeRepositories(reader, result, err)

	return result
}

func (sr *RepositoryReader) sendResult(result *Response) error {
	if done, err := sr.getErrors(result); done {
		return err
	}

	for _, edge := range result.Data.Search.Edges {
		node := edge.Node
		sr.output <- &node
	}

	return nil
}

func (sr *RepositoryReader) decodeRepositories(reader io.Reader, result *Response, err error) {
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

func (sr *RepositoryReader) repositoryQueryReader(query string) (io.ReadCloser, error) {
	request, _ := http.NewRequest("POST", "", strings.NewReader(query))
	response, err := sr.client.Do(request)

	if nil != err {
		return nil, err
	}

	return response.Body, nil
}

func (sr *RepositoryReader) buildQl(limit int, cursor string) string {
	if cursor != "" {
		cursor = fmt.Sprintf(", after: \\\"%s\\\"", cursor)
	}

	return fmt.Sprintf(repoSearchQuery, sr.query, limit, cursor)
}

func (sr *RepositoryReader) getErrors(result *Response) (bool, error) {
	if result.Message != "" {
		return true, fmt.Errorf("%s: %w", result.Message, ErrRead)
	}

	if 0 < len(result.Errors) {
		err := sr.wrapErrors(result)
		return true, err
	}

	return false, nil
}

func (sr *RepositoryReader) wrapErrors(result *Response) error {
	err := fmt.Errorf("api error: %w", ErrRead)

	for _, resultErr := range result.Errors {
		err = fmt.Errorf("%s - %s: %w", resultErr.Type, resultErr.Message, err)
	}

	return err
}

func NewRepositoryReader(query string, total int, output chan *Repository, client finderhttp.Client) *RepositoryReader {
	return &RepositoryReader{query: query, total: total, output: output, client: client, pageSize: 100}
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
