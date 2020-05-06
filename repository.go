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
		Path       string
		Keywords   []string
		Qualifiers []RepositorySearchQualifier
		Sort       RepositorySearchSort
		PerPage    int
	}
)

type Repository struct {
	Name     string
	FullName string
}
