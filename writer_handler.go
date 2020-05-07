package github_tool_finder

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type WriterHandler struct {
	input  chan *Repository
	closer io.Closer
	writer *csv.Writer
}

func (this WriterHandler) Handle() {

	for repository := range this.input {
		this.writeRepository(repository)
	}

	this.writer.Flush()
	this.closer.Close()
}

func (this WriterHandler) writeRepository(repository *Repository) {
	this.writeValues(
		repository.Name,
		repository.NameWithOwner,
		repository.Owner.Login,
		repository.Description,
		repository.Url,
		fmt.Sprintf("%d", repository.ForkCount),
		fmt.Sprintf("%d", repository.Stargazers.TotalCount),
		fmt.Sprintf("%d", repository.Watchers.TotalCount),
		repository.HomepageUrl,
		repository.LicenseInfo.Name,
		fmt.Sprintf("%d", repository.MentionableUsers.TotalCount),
		repository.MirrorUrl,
		strconv.FormatBool(repository.IsMirror),
		repository.PrimaryLanguage.Name,
		repository.Parent.Name,
		repository.CreatedAt.String(),
		repository.UpdatedAt.String(),
	)
}

func (this *WriterHandler) writeValues(values ...string) {
	this.writer.Write(values)
}

func NewWriterHandler(input chan *Repository, output io.WriteCloser) *WriterHandler {

	this := &WriterHandler{
		input:  input,
		closer: output,
		writer: csv.NewWriter(output),
	}

	this.writeValues("Name", "NameWithOwner", "Owner", "Description", "Url", "ForkCount", "Stargazers", "Watchers", "HomepageUrl", "LicenseInfo", "MentionableUsers", "MirrorUrl", "IsMirror", "PrimaryLanguage", "Parent", "CreatedAt", "UpdatedAt")

	return this
}
