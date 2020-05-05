package github_tool_finder

import (
	"bytes"
	"net/http"
	"strconv"
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
	searchConfig   *RepositorySearchConfig
}

func (this *SearchReaderFixture) Setup() {
	this.output = make(chan *Repository)
	this.fakeHttpClient = &FakeHTTPClient{}
	this.searchConfig = &RepositorySearchConfig{
		Url:      "https://api.github.com/search/repositories",
		Keywords: []string{"test1", "test2"},
		Qualifiers: []RepositorySearchQualifier{
			{Qualifier: "in", Value: "readme"},
			{Qualifier: "user", Value: "test"},
		},
		Sort:    RepositorySearchSort{},
		PerPage: 200,
	}
	this.searchReader = NewSearchReader(this.searchConfig, this.output, this.fakeHttpClient)
}

func (this *SearchReaderFixture) TestRead() {
	this.searchReader.Handle()
	//request := this.fakeHttpClient.request
	//this.So(request.URL, should.Equal, "https://api.github.com/search/repositories?q=test1+test2+in:readme+user:test&per_page=200")

	this.assertQueryString("q", "test1+test2+in:readme+user:test")
	this.assertQueryString("per_page", strconv.Itoa(this.searchConfig.PerPage))
}

func (this *SearchReaderFixture) assertQueryString(key, expected string) {
	query := this.fakeHttpClient.request.URL.Query()
	this.So(expected, should.Equal, query.Get(key))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
