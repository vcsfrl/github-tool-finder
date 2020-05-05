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
	query := make(url.Values)
	qualifiers := []string{}
	for _, qualifier := range this.filter.Qualifiers {
		val := fmt.Sprintf("%s:%s", qualifier.Qualifier, qualifier.Value)
		qualifiers = append(qualifiers, val)
	}

	keywords := strings.Join(append(this.filter.Keywords, qualifiers...), "+")

	query.Set("q", keywords)
	query.Set("per_page", strconv.Itoa(this.filter.PerPage))

	request, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", this.filter.Url, query.Encode()), nil)

	this.client.Do(request)
}

func NewSearchReader(filter *RepositorySearchConfig, output chan *Repository, client HTTPClient) *SearchReader {
	return &SearchReader{filter: filter, output: output, client: client}
}
