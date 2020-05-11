# github-tool-finder
Find tools on Github.

Search the Github api based on a query. Results are sent to STDOUT in CSV format.

### Install
 - git clone git@github.com:vcsfrl/github-tool-finder.git
 - `cd github-tool-finder/`
 - `make install`

### Usage
`./bin/search [query] [total]`
 - query: for details see the search section on https://developer.github.com/v4/query/
 - total: maximum number of results to fetch

ENV Variables:
 - GH_TOKEN - oAuth access token from Github.

### Tests
 - `cd /project/path`
 - Run tests: `make test`
 - Run coverage: `make cover`
 - Run coverage html: `make cover-html`

### Examples
 - `./bin/search "orm language:php sort:stars-desc" 50 > /path/to/result.csv`
 - `GH_TOKEN=github_access_token ./bin/search "orm language:php sort:stars-desc" 50 > /path/to/result.csv`