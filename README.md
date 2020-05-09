# github-tool-finder
Find tools on Github.

### Install

`cd /project/path`

`make install`

### Usage

`./bin/search [query] [total]`


##### ENV Variables
GH_TOKEN - oAuth access token from github.

### Run tests

`make test`

### Example
`GH_TOKEN=github_access_token ./bin/search "orm language:php sort:stars-desc" 50 > /path/to/result.csv`