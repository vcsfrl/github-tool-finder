package github_tool_finder

import (
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

	inner  *FakeHTTPClient
	client *AuthenticationClientV4
}

func (this *AuthenticationClientFixture) Setup() {
	this.inner = &FakeHTTPClient{}
	this.client = NewAuthenticationClientV4(this.inner, "authtoken")
}

func (this *AuthenticationClientFixture) TestResponseFromInnerClientReturned() {
	this.inner.response = &http.Response{StatusCode: http.StatusTeapot}
	this.inner.err = errors.New("HTTP Error")
	request := httptest.NewRequest("GET", "/path", nil)
	response, err := this.client.Do(request)
	this.So(response.StatusCode, should.Equal, http.StatusTeapot)
	this.So(err.Error(), should.Equal, "HTTP Error")
}

func (this *AuthenticationClientFixture) TestProvidedInformationAddedBeforeRequestSent() {
	request := httptest.NewRequest("GET", "/path?existingKey=existingValue", nil)

	this.client.Do(request)
	this.assertRequestConnectionInformation()
	this.assertQueryStringIncludesAuthentication()
	this.assertQueryStringValue("existingKey", "existingValue")
}

func (this *AuthenticationClientFixture) assertRequestConnectionInformation() {
	this.So(this.inner.request.URL.Scheme, should.Equal, "https")
	this.So(this.inner.request.URL.Host, should.Equal, "api.github.com")
	this.So(this.inner.request.Host, should.Equal, "api.github.com")
	this.So(this.inner.request.URL.Path, should.Equal, "graphql")
}

func (this *AuthenticationClientFixture) TestMissingToken() {
	request := httptest.NewRequest("GET", "/path?existingKey=existingValue", nil)
	this.client.authToken = ""
	this.client.Do(request)
	this.assertRequestConnectionInformation()
	this.So(this.inner.request.Header.Get("Authorization"), should.Equal, "")

}

func (this *AuthenticationClientFixture) assertQueryStringIncludesAuthentication() {
	this.So(this.inner.request.Header.Get("Authorization"), should.Equal, "bearer authtoken")

}

func (this *AuthenticationClientFixture) assertQueryStringValue(key string, expectedValue string) {
	this.So(this.inner.request.URL.Query().Get(key), should.Equal, expectedValue)
}
