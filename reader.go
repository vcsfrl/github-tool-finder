package github_tool_finder

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type SearchReader struct {
	filter *RepositorySearchConfig
	output chan *Repository
	client HTTPClient
}

func (this *SearchReader) Handle() {
	request, _ := http.NewRequest("GET", this.buildUrl(), nil)
	this.client.Do(request)
}

func (this *SearchReader) buildUrl() string {
	return fmt.Sprintf("%s?%s", this.filter.Path, this.buildQuery().Encode())
}

func (this *SearchReader) buildQuery() url.Values {
	query, parts := make(url.Values), []string{}
	for _, qualifier := range this.filter.Qualifiers {
		parts = append(parts, fmt.Sprintf("%s:%s", qualifier.Qualifier, qualifier.Value))
	}

	query.Set("q", strings.Join(append(this.filter.Keywords, parts...), "+"))
	query.Set("per_page", strconv.Itoa(this.filter.PerPage))

	return query
}

func NewSearchReader(filter *RepositorySearchConfig, output chan *Repository, client HTTPClient) *SearchReader {
	return &SearchReader{filter: filter, output: output, client: client}
}
