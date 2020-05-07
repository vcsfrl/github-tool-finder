package github_tool_finder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type SearchReader struct {
	client HTTPClient
	query  string
	number int
	output chan *Repository
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

		log.Println(edge.Node)
		this.output <- &edge.Node
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

	res, _ := ioutil.ReadAll(reader)
	log.Fatal(string(res))
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(result)

	if nil != err {
		result.ErrorMessage = err.Error()
	}
}

func (this *SearchReader) repositoryReader() (io.ReadCloser, error) {
	request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(this.buildQl())))
	if nil != err {
		return nil, err
	}
	response, err := this.client.Do(request)

	log.Fatal(request.URL.String())
	if nil != err {
		return nil, err
	}

	return response.Body, nil
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
