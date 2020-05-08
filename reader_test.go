package github_tool_finder

import (
	"bytes"
	"errors"
	"fmt"
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
	this.searchReader.pageSize = 1
}

func (this *SearchReaderFixture) TestReadResponse() {
	this.fakeClient.Configure(responseBody, 200, nil)
	this.searchReader.Handle()

	body, _ := ioutil.ReadAll(this.fakeClient.request.Body)
	this.So(string(body), should.Equal, grapqlQuery1Result)
	this.So(<-this.output, should.Resemble, getResponseRepository(1))
	this.So(this.fakeClient.responseBody.closed, should.Equal, 1)
	this.So(this.fakeClient.callNr, should.Equal, 1)
}

func (this *SearchReaderFixture) SkipTestPaginatedRead() {
	this.searchReader.nrRepos = 2

	this.fakeClient.Configure(responseBody, 200, nil)
	this.searchReader.Handle()

	body, _ := ioutil.ReadAll(this.fakeClient.request.Body)
	this.So(string(body), should.Equal, grapqlQuery2Result)
	this.So(<-this.output, should.Resemble, getResponseRepository(2))
	this.So(this.fakeClient.responseBody.closed, should.Equal, 1)
	this.So(this.fakeClient.callNr, should.Equal, 2)
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
	request            *http.Request
	response           *http.Response
	responseBody       *SearchReaderBuffer
	responseContent    []string
	responseStatusCode int
	err                error
	callNr             int
}

func (this *FakeHTTPClient) Configure(responseContent []string, statusCode int, err error) {
	if err == nil {
		this.responseContent = responseContent
		this.responseStatusCode = statusCode
	}
	this.err = err
}

func (this *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {

	this.request = request

	if nil != this.err {
		return nil, this.err
	}

	this.responseBody = NewSearchReadBuffer(this.responseContent[this.callNr])
	this.response = &http.Response{
		Body:       this.responseBody,
		StatusCode: this.responseStatusCode,
	}
	this.callNr++

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

func getResponseRepository(index int) *Repository {

	created, _ := time.Parse(time.RFC3339, "2015-05-23T21:24:16Z")
	updated, _ := time.Parse(time.RFC3339, "2020-04-15T20:01:25Z")

	return &Repository{
		Description:   fmt.Sprintf("%d Test description.", index),
		Name:          fmt.Sprintf("%dtestrepo", index),
		NameWithOwner: fmt.Sprintf("%dtestrepo/testrepo", index),
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

const grapqlQuery1Result = "{\"query\":\"query SearchRepositories {\\n  search(query: \\\"test:test test\\\", type: REPOSITORY, first:1){\\n    repositoryCount\\n    edges {\\n      cursor \\n      node {\\n\\t\\t\\t\\t... on Repository {\\n          description\\n          name\\n          nameWithOwner\\n          url\\n          owner {\\n            login\\n          }\\n          forkCount\\n          stargazers {\\n            totalCount\\n          }\\n          watchers {\\n            totalCount\\n          }\\n          homepageUrl\\n          licenseInfo {\\n            name\\n          }\\n          mentionableUsers {\\n            totalCount\\n          }\\n          mirrorUrl\\n          isMirror\\n          primaryLanguage {\\n            name\\n          }\\n          parent {\\n            name\\n          }\\n          createdAt\\n          updatedAt\\n        }\\n      }\\n    }\\n  }\\n}\\n\",\"variables\":{}}"
const grapqlQuery2Result = "{\"query\":\"query SearchRepositories {\\n  search(query: \\\"test:test test\\\", type: REPOSITORY, first:1, after: \\\"aaa\\\"){\\n    repositoryCount\\n    edges {\\n      cursor \\n      node {\\n\\t\\t\\t\\t... on Repository {\\n          description\\n          name\\n          nameWithOwner\\n          url\\n          owner {\\n            login\\n          }\\n          forkCount\\n          stargazers {\\n            totalCount\\n          }\\n          watchers {\\n            totalCount\\n          }\\n          homepageUrl\\n          licenseInfo {\\n            name\\n          }\\n          mentionableUsers {\\n            totalCount\\n          }\\n          mirrorUrl\\n          isMirror\\n          primaryLanguage {\\n            name\\n          }\\n          parent {\\n            name\\n          }\\n          createdAt\\n          updatedAt\\n        }\\n      }\\n    }\\n  }\\n}\\n\",\"variables\":{}}"

var responseBody = []string{
	`{
    "data": {
        "search": {
            "repositoryCount": 128,
            "edges": [
                {
					"cursor": "aaa",
                    "node": {
                        "description": "1 Test description.",
                        "name": "1testrepo",
                        "nameWithOwner": "1testrepo/testrepo",
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
}`,
	`{
    "data": {
        "search": {
            "repositoryCount": 128,
            "edges": [
                {
					"cursor": "bbb",
                    "node": {
                        "description": "2 Test description.",
                        "name": "2testrepo",
                        "nameWithOwner": "2testrepo/testrepo",
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
}`,
}

var responseError = []string{`{
    "message": "Bad credentials",
    "documentation_url": "https://developer.github.com/v4"
}`,
}
