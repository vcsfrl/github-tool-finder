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

	handler *WriterHandler
	input   chan *Repository
	buffer  *ReadWriteSpyBuffer
	writer  *csv.Writer
}

func (this *WriterHandlerFixture) Setup() {
	this.buffer = NewReadWriteSpyBuffer("")
	this.input = make(chan *Repository, 10)
	this.handler = NewWriterHandler(this.input, this.buffer)
}

func (this *WriterHandlerFixture) TestOutputClosed() {
	close(this.input)
	this.handler.Handle()
	this.So(this.buffer.closed, should.Equal, 1)
}

func (this *WriterHandlerFixture) TestHeaderMatchesRecord() {
	this.input <- this.createRepository(1)
	close(this.input)
	this.handler.Handle()
	this.assertHeaderMatchesRecord()
}

func (this *WriterHandlerFixture) assertHeaderMatchesRecord() {
	lines := this.outputLines()
	header := lines[0]
	record := lines[1]

	this.So(header, should.Equal, "Name,NameWithOwner,Owner,Description,Url,ForkCount,Stargazers,Watchers,HomepageUrl,LicenseInfo,MentionableUsers,MirrorUrl,IsMirror,PrimaryLanguage,Parent,CreatedAt,UpdatedAt")
	this.So(record, should.Equal, "Name1,NameWithOwner1,Owner1,Description1,Url1,2,3,4,HomepageUrl1,LicenseInfo1,5,MirrorUrl1,false,PrimaryLanguage1,Parent1,2020-04-15 20:01:25 +0000 UTC,2020-05-15 20:01:25 +0000 UTC")

}

func (this *WriterHandlerFixture) TestAllRepositoriesWritten() {
	this.sendEnvelopes(2)
	this.handler.Handle()

	if lines := this.outputLines(); this.So(lines, should.HaveLength, 3) {
		this.So(lines[1], should.Equal, "Name1,NameWithOwner1,Owner1,Description1,Url1,2,3,4,HomepageUrl1,LicenseInfo1,5,MirrorUrl1,false,PrimaryLanguage1,Parent1,2020-04-15 20:01:25 +0000 UTC,2020-05-15 20:01:25 +0000 UTC")
		this.So(lines[2], should.Equal, "Name2,NameWithOwner2,Owner2,Description2,Url2,3,4,5,HomepageUrl2,LicenseInfo2,6,MirrorUrl2,false,PrimaryLanguage2,Parent2,2020-04-15 20:01:25 +0000 UTC,2020-05-15 20:01:25 +0000 UTC")
	}
}

func (this *WriterHandlerFixture) sendEnvelopes(count int) {
	for i := 1; i < count+1; i++ {
		this.input <- this.createRepository(int64(i))
	}
	close(this.input)
}

func (this *WriterHandlerFixture) createRepository(index int64) *Repository {

	created, _ := time.Parse(time.RFC3339, "2020-04-15T20:01:25Z")
	updated, _ := time.Parse(time.RFC3339, "2020-05-15T20:01:25Z")
	return &Repository{
		Description:   fmt.Sprintf("Description%d", index),
		Name:          fmt.Sprintf("Name%d", index),
		NameWithOwner: fmt.Sprintf("NameWithOwner%d", index),
		Url:           fmt.Sprintf("Url%d", index),
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
		HomepageUrl: fmt.Sprintf("HomepageUrl%d", index),
		LicenseInfo: struct {
			Name string `json:"name"`
		}{Name: fmt.Sprintf("LicenseInfo%d", index)},
		MentionableUsers: struct {
			TotalCount int64 `json:"totalCount"`
		}{TotalCount: index + 4},
		MirrorUrl: fmt.Sprintf("MirrorUrl%d", index),
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

func (this *WriterHandlerFixture) outputLines() []string {
	outputFile := strings.TrimSpace(this.buffer.String())
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
