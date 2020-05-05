package github_tool_finder

type (
	RepositorySearchSort struct {
		Field string
		Asc   bool
	}
	RepositorySearchQualifier struct {
		Qualifier string
		Value     string
	}
	RepositorySearchConfig struct {
		Url        string
		Keywords   []string
		Qualifiers []RepositorySearchQualifier
		Sort       RepositorySearchSort
		PerPage    int
	}
)

type Repository struct {
}
