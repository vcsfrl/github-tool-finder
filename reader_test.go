package github_tool_finder

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/smartystreets/assertions/should"

	"github.com/smartystreets/gunit"
)

func TestSearchReaderFixture(t *testing.T) {
	gunit.Run(new(SearchReaderFixture), t)
}

type SearchReaderFixture struct {
	*gunit.Fixture

	fakeHttpClient *FakeHTTPClient
	output         chan *Repository
	searchReader   *SearchReader
}

func (this *SearchReaderFixture) Setup() {
	this.output = make(chan *Repository)
	this.fakeHttpClient = &FakeHTTPClient{}

	this.searchReader = NewSearchReader("test:test test", 10, this.output, this.fakeHttpClient)
}

func (this *SearchReaderFixture) TestBuildQuery() {
	this.searchReader.Handle()
	request := this.fakeHttpClient.request
	body, _ := ioutil.ReadAll(request.Body)
	this.So(string(body), should.Equal, grapqlQueryResult)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const grapqlQueryResult = `{
  search(query: "test:test test", type: REPOSITORY, first: 10) {
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

type FakeHTTPClient struct {
	request      *http.Request
	response     *http.Response
	responseBody *SerchReaderBuffer
	err          error
}

func (this *FakeHTTPClient) Configure(responseText string, statusCode int, err error) {
	if err == nil {
		this.responseBody = NewSearchReadBuffer(responseText)
		this.response = &http.Response{
			Body:       this.responseBody,
			StatusCode: statusCode,
		}
	}
	this.err = err
}

func (this *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	this.request = request

	return this.response, this.err
}

//////////

type SerchReaderBuffer struct {
	*bytes.Buffer
	closed int
}

func NewSearchReadBuffer(value string) *SerchReaderBuffer {
	return &SerchReaderBuffer{
		Buffer: bytes.NewBufferString(value),
	}
}

func (this *SerchReaderBuffer) Close() error {
	this.closed++
	this.Buffer.Reset()

	return nil
}
