package github_tool_finder

import (
	"fmt"
	"net/http"
)

func NewAuthenticationClient(inner HTTPClient, scheme string, hostname string, authToken string) *AuthenticationClient {
	return &AuthenticationClient{
		inner:     inner,
		scheme:    scheme,
		hostname:  hostname,
		authToken: authToken,
	}
}

type AuthenticationClient struct {
	scheme    string
	hostname  string
	inner     HTTPClient
	authToken string
}

func (this AuthenticationClient) Do(request *http.Request) (*http.Response, error) {
	request.URL.Scheme = this.scheme
	request.URL.Host = this.hostname
	if "" != this.authToken {
		request.Header.Set("Authorization", fmt.Sprintf("bearer %s", this.authToken))
	}
	request.Host = this.hostname
	return this.inner.Do(request)
}
