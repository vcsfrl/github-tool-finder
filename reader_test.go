package github_tool_finder

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

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

	this.searchReader = NewSearchReader("test:test test", 1, this.output, this.fakeHttpClient)
}

func (this *SearchReaderFixture) TestBuildQuery() {
	this.searchReader.Handle()
	request := this.fakeHttpClient.request
	body, _ := ioutil.ReadAll(request.Body)
	this.So(string(body), should.Equal, grapqlQueryResult)
}

func (this *SearchReaderFixture) TestReadResponse() {
	this.fakeHttpClient.Configure(responseBody, 200, nil)
	this.searchReader.Handle()
	this.So(<-this.output, should.Equal, getResponseRepository())
	this.So(this.fakeHttpClient.responseBody.closed, should.Equal, 1)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FakeHTTPClient struct {
	request      *http.Request
	response     *http.Response
	responseBody *SearchReaderBuffer
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

type SearchReaderBuffer struct {
	*bytes.Buffer
	closed int
}

//////////

func NewSearchReadBuffer(value string) *SearchReaderBuffer {
	return &SearchReaderBuffer{
		Buffer: bytes.NewBufferString(value),
	}
}

func (this *SearchReaderBuffer) Close() error {
	this.closed++
	this.Buffer.Reset()

	return nil
}

func getResponseRepository() Repository {
	return Repository{
		Description:      "Test description.",
		Name:             "testrepo",
		NameWithOwner:    "testrepo/",
		Url:              "",
		Owner:            "",
		ForkCount:        0,
		Stargazers:       0,
		Watchers:         0,
		HomepageUrl:      "",
		LicenseInfo:      "",
		MentionableUsers: 0,
		MirrorUrl:        "",
		IsMirror:         false,
		PrimaryLanguage:  "",
		Parent:           "",
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
	}
}

//////////

const grapqlQueryResult = `{
  search(query: "test:test test", type: REPOSITORY, first: 1) {
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

const responseBody = `{
    "data": {
        "search": {
            "repositoryCount": 128,
            "edges": [
                {
                    "node": {
                        "description": "Test description.",
                        "name": "testrepo",
                        "nameWithOwner": "testrepo/testrepo",
                        "url": "https://github.com/testrepo/testrepo",
                        "owner": {
                            "login": "testrepo"
                        },
                        "forkCount": 10,
                        "stargazers": {
                            "totalCount": 10
                        },
                        "watchers": {
                            "totalCount": 10
                        },
                        "homepageUrl": null,
                        "licenseInfo": {
                            "name": "Test license"
                        },
                        "mentionableUsers": {
                            "totalCount": 10
                        },
                        "mirrorUrl": null,
                        "isMirror": false,
                        "primaryLanguage": {
                            "name": "Go"
                        },
                        "parent": null,
                        "createdAt": "2015-05-23T21:24:16Z",
                        "updatedAt": "2020-04-15T20:01:25Z"
                    }
                }
            ]
        }
    }
}`
