package search

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type CsvWriter struct {
	input  chan *Repository
	closer io.Closer
	writer *csv.Writer
}

func (cw *CsvWriter) Handle() {
	for repository := range cw.input {
		cw.writeRepository(repository)
	}

	cw.writer.Flush()
	cw.closer.Close()
}

func (cw *CsvWriter) writeRepository(repository *Repository) {
	cw.writeValues(
		repository.Name,
		repository.NameWithOwner,
		repository.Owner.Login,
		repository.Description,
		repository.URL,
		fmt.Sprintf("%d", repository.ForkCount),
		fmt.Sprintf("%d", repository.Stargazers.TotalCount),
		fmt.Sprintf("%d", repository.Watchers.TotalCount),
		repository.HomepageURL,
		repository.LicenseInfo.Name,
		fmt.Sprintf("%d", repository.MentionableUsers.TotalCount),
		repository.MirrorURL,
		strconv.FormatBool(repository.IsMirror),
		repository.PrimaryLanguage.Name,
		repository.Parent.Name,
		repository.CreatedAt.String(),
		repository.UpdatedAt.String(),
	)
}

func (cw *CsvWriter) writeValues(values ...string) {
	cw.writer.Write(values)
}

func NewCsvWriter(input chan *Repository, output io.WriteCloser) *CsvWriter {
	this := &CsvWriter{
		input:  input,
		closer: output,
		writer: csv.NewWriter(output),
	}

	this.writeValues("Name", "NameWithOwner", "Owner", "Description", "URL", "ForkCount", "Stargazers", "Watchers", "HomepageURL", "LicenseInfo", "MentionableUsers", "MirrorURL", "IsMirror", "PrimaryLanguage", "Parent", "CreatedAt", "UpdatedAt")

	return this
}
