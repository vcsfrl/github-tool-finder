package github_tool_finder

import (
	"bytes"
	"errors"
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

	fakeClient   *FakeHTTPClient
	output       chan *Repository
	searchReader *SearchReader
}

func (this *SearchReaderFixture) Setup() {
	this.output = make(chan *Repository, 10)
	this.fakeClient = &FakeHTTPClient{}

	this.searchReader = NewSearchReader("test:test test", 1, this.output, this.fakeClient)
}

func (this *SearchReaderFixture) TestBuildQuery() {
	this.fakeClient.Configure(responseBody, 200, nil)
	this.searchReader.Handle()
	request := this.fakeClient.request
	body, _ := ioutil.ReadAll(request.Body)
	this.So(string(body), should.Equal, grapqlQueryResult)
}

func (this *SearchReaderFixture) TestReadResponse() {
	this.fakeClient.Configure(responseBody, 200, nil)
	this.searchReader.Handle()
	this.So(<-this.output, should.Resemble, getResponseRepository())
	this.So(this.fakeClient.responseBody.closed, should.Equal, 1)
}

func (this *SearchReaderFixture) TestReadReadError() {
	this.fakeClient.Configure(responseError, 401, nil)
	err := this.searchReader.Handle()
	this.So(err, should.BeError)
	this.So(err.Error(), should.Equal, "Bad credentials")
	this.So(this.fakeClient.responseBody.closed, should.Equal, 1)
}

func (this *SearchReaderFixture) TestReadError() {
	this.fakeClient.Configure(responseError, 401, errors.New("test error"))
	err := this.searchReader.Handle()
	this.So(err, should.BeError)
	this.So(err.Error(), should.Equal, "test error")
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

func getResponseRepository() *Repository {

	created, _ := time.Parse(time.RFC3339, "2015-05-23T21:24:16Z")
	updated, _ := time.Parse(time.RFC3339, "2020-04-15T20:01:25Z")

	return &Repository{
		Description:   "Test description.",
		Name:          "testrepo",
		NameWithOwner: "testrepo/testrepo",
		Url:           "https://github.com/testrepo/testrepo",
		Owner: struct {
			Login string `json:"login"`
		}{Login: "testrepo"},
		ForkCount: 10,
		Stargazers: struct {
			TotalCount int64 `json:"totalCount"`
		}{TotalCount: 10},
		Watchers: struct {
			TotalCount int64 `json:"totalCount"`
		}{TotalCount: 10},
		HomepageUrl: "testhomepage",
		LicenseInfo: struct {
			Name string `json:"name"`
		}{Name: "testlicense"},
		MentionableUsers: struct {
			TotalCount int64 `json:"totalCount"`
		}{TotalCount: 10},
		MirrorUrl: "testmirror",
		IsMirror:  true,
		PrimaryLanguage: struct {
			Name string `json:"name"`
		}{Name: "Go"},
		Parent: struct {
			Name string `json:"name"`
		}{Name: "testparent"},
		CreatedAt: created,
		UpdatedAt: updated,
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
                        "homepageUrl": "testhomepage",
                        "licenseInfo": {
                            "name": "testlicense"
                        },
                        "mentionableUsers": {
                            "totalCount": 10
                        },
                        "mirrorUrl": "testmirror",
                        "isMirror": true,
                        "primaryLanguage": {
                            "name": "Go"
                        },
                        "parent": {
                            "name": "testparent"
						},
                        "createdAt": "2015-05-23T21:24:16Z",
                        "updatedAt": "2020-04-15T20:01:25Z"
                    }
                }
            ]
        }
    }
}`

const responseError = `{
    "message": "Bad credentials",
    "documentation_url": "https://developer.github.com/v4"
}`
