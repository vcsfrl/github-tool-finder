package github_tool_finder

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"

	"github.com/smartystreets/gunit"
)

func TestWriterHandlerFixture(t *testing.T) {
	gunit.Run(new(WriterHandlerFixture), t)
}

type WriterHandlerFixture struct {
	*gunit.Fixture

	handler *CsvWriter
	input   chan *Repository
	buffer  *ReadWriteSpyBuffer
	writer  *csv.Writer
}

func (whf *WriterHandlerFixture) Setup() {
	whf.buffer = NewReadWriteSpyBuffer("")
	whf.input = make(chan *Repository, 10)
	whf.handler = NewCsvWriter(whf.input, whf.buffer)
}

func (whf *WriterHandlerFixture) TestOutputClosed() {
	close(whf.input)
	whf.handler.Handle()
	whf.So(whf.buffer.closed, should.Equal, 1)
}

func (whf *WriterHandlerFixture) TestHeaderMatchesRecord() {
	whf.input <- whf.createRepository(1)
	close(whf.input)
	whf.handler.Handle()
	whf.assertHeaderMatchesRecord()
}

func (whf *WriterHandlerFixture) assertHeaderMatchesRecord() {
	lines := whf.outputLines()
	header := lines[0]
	record := lines[1]

	whf.So(header, should.Equal, "Name,NameWithOwner,Owner,Description,URL,ForkCount,Stargazers,Watchers,HomepageURL,LicenseInfo,MentionableUsers,MirrorURL,IsMirror,PrimaryLanguage,Parent,CreatedAt,UpdatedAt")
	whf.So(record, should.Equal, "Name1,NameWithOwner1,Owner1,Description1,URL1,2,3,4,HomepageURL1,LicenseInfo1,5,MirrorURL1,false,PrimaryLanguage1,Parent1,2020-04-15 20:01:25 +0000 UTC,2020-05-15 20:01:25 +0000 UTC")

}

func (whf *WriterHandlerFixture) TestAllRepositoriesWritten() {
	whf.sendEnvelopes(2)
	whf.handler.Handle()

	if lines := whf.outputLines(); whf.So(lines, should.HaveLength, 3) {
		whf.So(lines[1], should.Equal, "Name1,NameWithOwner1,Owner1,Description1,URL1,2,3,4,HomepageURL1,LicenseInfo1,5,MirrorURL1,false,PrimaryLanguage1,Parent1,2020-04-15 20:01:25 +0000 UTC,2020-05-15 20:01:25 +0000 UTC")
		whf.So(lines[2], should.Equal, "Name2,NameWithOwner2,Owner2,Description2,URL2,3,4,5,HomepageURL2,LicenseInfo2,6,MirrorURL2,false,PrimaryLanguage2,Parent2,2020-04-15 20:01:25 +0000 UTC,2020-05-15 20:01:25 +0000 UTC")
	}
}

func (whf *WriterHandlerFixture) sendEnvelopes(count int) {
	for i := 1; i < count+1; i++ {
		whf.input <- whf.createRepository(int64(i))
	}
	close(whf.input)
}

func (whf *WriterHandlerFixture) createRepository(index int64) *Repository {

	created, _ := time.Parse(time.RFC3339, "2020-04-15T20:01:25Z")
	updated, _ := time.Parse(time.RFC3339, "2020-05-15T20:01:25Z")
	return &Repository{
		Description:   fmt.Sprintf("Description%d", index),
		Name:          fmt.Sprintf("Name%d", index),
		NameWithOwner: fmt.Sprintf("NameWithOwner%d", index),
		URL:           fmt.Sprintf("URL%d", index),
		Owner: struct {
			Login string `json:"login"`
		}{Login: fmt.Sprintf("Owner%d", index)},
		ForkCount: index + 1,
		Stargazers: struct {
			TotalCount int64 `json:"totalCount"`
		}{TotalCount: index + 2},
		Watchers: struct {
			TotalCount int64 `json:"totalCount"`
		}{TotalCount: index + 3},
		HomepageURL: fmt.Sprintf("HomepageURL%d", index),
		LicenseInfo: struct {
			Name string `json:"name"`
		}{Name: fmt.Sprintf("LicenseInfo%d", index)},
		MentionableUsers: struct {
			TotalCount int64 `json:"totalCount"`
		}{TotalCount: index + 4},
		MirrorURL: fmt.Sprintf("MirrorURL%d", index),
		IsMirror:  false,
		PrimaryLanguage: struct {
			Name string `json:"name"`
		}{Name: fmt.Sprintf("PrimaryLanguage%d", index)},
		Parent: struct {
			Name string `json:"name"`
		}{Name: fmt.Sprintf("Parent%d", index)},
		CreatedAt: created,
		UpdatedAt: updated,
	}
}

func (whf *WriterHandlerFixture) outputLines() []string {
	outputFile := strings.TrimSpace(whf.buffer.String())
	return strings.Split(outputFile, "\n")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ReadWriteSpyBuffer struct {
	*bytes.Buffer
	closed int
}

func NewReadWriteSpyBuffer(value string) *ReadWriteSpyBuffer {
	return &ReadWriteSpyBuffer{
		Buffer: bytes.NewBufferString(value),
	}
}

func (this *ReadWriteSpyBuffer) Close() error {
	this.closed++

	return nil
}
