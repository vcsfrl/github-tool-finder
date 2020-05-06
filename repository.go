package github_tool_finder

import "time"

type SearchResponse struct {
	Data struct {
		Search struct {
			RepositoryCount int `json:"repositoryCount"`
			Edges           []struct {
				Node Repository `json:"node"`
			} `json:"edges"`
		} `json:"search"`
	} `json:"data"`
}

type Repository struct {
	Description      string    `json:"description"`
	Name             string    `json:"name"`
	NameWithOwner    string    `json:"nameWithOwner"`
	Url              string    `json:"url"`
	Owner            string    `json:"owner.login"`
	ForkCount        int64     `json:"forkCount"`
	Stargazers       int64     `json:"stargazers.totalCount"`
	Watchers         int64     `json:"watchers.totalCount"`
	HomepageUrl      string    `json:"homepageUrl"`
	LicenseInfo      string    `json:"licenseInfo.name"`
	MentionableUsers int       `json:"mentionableUsers.totalCount"`
	MirrorUrl        string    `json:"mirrorUrl"`
	IsMirror         bool      `json:"isMirror"`
	PrimaryLanguage  string    `json:"primaryLanguage.name"`
	Parent           string    `json:"parent.name"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
