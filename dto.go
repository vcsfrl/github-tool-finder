package github_tool_finder

import "time"

type SearchResponse struct {
	Data struct {
		Search struct {
			RepositoryCount int `json:"repositoryCount"`
			Edges           []struct {
				Cursor string     `json:"cursor"`
				Node   Repository `json:"node"`
			} `json:"edges"`
		} `json:"search"`
	} `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Errors  []struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

type Repository struct {
	Description   string `json:"description"`
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
	URL           string `json:"url"`
	Owner         struct {
		Login string `json:"login"`
	} `json:"owner"`
	ForkCount  int64 `json:"forkCount"`
	Stargazers struct {
		TotalCount int64 `json:"totalCount"`
	} `json:"stargazers"`
	Watchers struct {
		TotalCount int64 `json:"totalCount"`
	} `json:"watchers"`
	HomepageURL string `json:"homepageUrl"`
	LicenseInfo struct {
		Name string `json:"name"`
	} `json:"licenseInfo"`
	MentionableUsers struct {
		TotalCount int64 `json:"totalCount"`
	} `json:"mentionableUsers"`
	MirrorURL       string `json:"mirrorUrl"`
	IsMirror        bool   `json:"isMirror"`
	PrimaryLanguage struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`
	Parent struct {
		Name string `json:"name"`
	} `json:"parent"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
