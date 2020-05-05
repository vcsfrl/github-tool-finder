package github_tool_finder

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestSearchReaderFixture(t *testing.T) {
	gunit.Run(new(SearchReaderFixture), t)
}

type SearchReaderFixture struct {
	*gunit.Fixture
}

func (this *SearchReaderFixture) Setup() {
}

func (this *SearchReaderFixture) TestVerifierReceivesInput() {

}
