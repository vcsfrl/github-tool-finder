package github_tool_finder

import (
	"fmt"
	"net/http"
)

func NewAuthenticationClientV4(inner HTTPClient, authToken string) *AuthenticationClientV4 {
	return &AuthenticationClientV4{
		inner:     inner,
		authToken: authToken,
	}
}

type AuthenticationClientV4 struct {
	inner     HTTPClient
	authToken string
}

func (ac *AuthenticationClientV4) Do(request *http.Request) (*http.Response, error) {
	request.URL.Scheme = "https"
	request.URL.Host = "api.github.com"
	request.Host = "api.github.com"
	request.URL.Path = "/graphql"
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	if "" != ac.authToken {
		request.Header.Add("Authorization", fmt.Sprintf("bearer %s", ac.authToken))
	}

	return ac.inner.Do(request)
}
