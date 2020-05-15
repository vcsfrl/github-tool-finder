package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"

	"github.com/smartystreets/gunit"
)

func TestAuthenticationClient(t *testing.T) {
	gunit.Run(new(AuthenticationClientFixture), t)
}

type AuthenticationClientFixture struct {
	*gunit.Fixture

	inner  *FakeSipleHTTPClient
	client *AuthenticationClientV4
}

func (acf *AuthenticationClientFixture) Setup() {
	acf.inner = &FakeSipleHTTPClient{}
	acf.client = NewAuthenticationClientV4(acf.inner, "authtoken")
}

func (acf *AuthenticationClientFixture) TestResponseFromInnerClientReturned() {
	acf.inner.response = &http.Response{StatusCode: http.StatusTeapot}
	acf.inner.err = errors.New("HTTP Error")
	request := httptest.NewRequest("GET", "/path", nil)
	response, err := acf.client.Do(request)
	acf.So(response.StatusCode, should.Equal, http.StatusTeapot)
	acf.So(err.Error(), should.Equal, "HTTP Error")
}

func (acf *AuthenticationClientFixture) TestProvidedInformationAddedBeforeRequestSent() {
	request := httptest.NewRequest("GET", "/path?existingKey=existingValue", nil)

	acf.client.Do(request)
	acf.assertRequestConnectionInformation()
	acf.assertQueryStringIncludesAuthentication()
	acf.assertQueryStringValue("existingKey", "existingValue")
}

func (acf *AuthenticationClientFixture) assertRequestConnectionInformation() {
	acf.So(acf.inner.request.URL.Scheme, should.Equal, "https")
	acf.So(acf.inner.request.URL.Host, should.Equal, "api.github.com")
	acf.So(acf.inner.request.Host, should.Equal, "api.github.com")
	acf.So(acf.inner.request.URL.Path, should.Equal, "/graphql")
}

func (acf *AuthenticationClientFixture) TestMissingToken() {
	request := httptest.NewRequest("GET", "/path?existingKey=existingValue", nil)
	acf.client.authToken = ""
	acf.client.Do(request)
	acf.assertRequestConnectionInformation()
	acf.So(acf.inner.request.Header.Get("Authorization"), should.Equal, "")
}

func (acf *AuthenticationClientFixture) assertQueryStringIncludesAuthentication() {
	acf.So(acf.inner.request.Header.Get("Authorization"), should.Equal, "bearer authtoken")
}

func (acf *AuthenticationClientFixture) assertQueryStringValue(key string, expectedValue string) {
	acf.So(acf.inner.request.URL.Query().Get(key), should.Equal, expectedValue)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FakeSipleHTTPClient struct {
	request      *http.Request
	response     *http.Response
	responseBody *SearchReaderBuffer
	err          error
}

func (fc *FakeSipleHTTPClient) Configure(responseText string, statusCode int, err error) {
	if err == nil {
		fc.responseBody = NewSearchReadBuffer(responseText)
		fc.response = &http.Response{
			Body:       fc.responseBody,
			StatusCode: statusCode,
		}
	}

	fc.err = err
}

func (fc *FakeSipleHTTPClient) Do(request *http.Request) (*http.Response, error) {
	fc.request = request

	return fc.response, fc.err
}

type SearchReaderBuffer struct {
	*bytes.Buffer
	closed int
}

func NewSearchReadBuffer(value string) *SearchReaderBuffer {
	return &SearchReaderBuffer{
		Buffer: bytes.NewBufferString(value),
	}
}

func (sr *SearchReaderBuffer) Close() error {
	sr.closed++
	sr.Buffer.Reset()

	return nil
}
