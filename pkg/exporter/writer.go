package exporter

import (
	"fmt"
	"io"

	"github.com/daemonp/gogpt/pkg/fileutils"
)

type Writer struct {
	output io.Writer
}

func NewWriter(output io.Writer) *Writer {
	return &Writer{output: output}
}

func (w *Writer) Write(content string) {
	fmt.Fprint(w.output, content)
}

func (w *Writer) WriteFileContents(files []FileInfo) error {
	for _, file := range files {
		w.writeFileContent(file)
	}
	return nil
}

func (w *Writer) writeFileContent(file FileInfo) {
	fmt.Fprintf(w.output, "// File: %s\n", file.Path)
	if file.Excluded {
		fmt.Fprintf(w.output, "%s\n\n", file.Content)
		return
	}
	fmt.Fprintf(w.output, "```%s\n", fileutils.GetFileExtension(file.Path))
	fmt.Fprintf(w.output, "%s\n", file.Content)
	fmt.Fprintf(w.output, "```\n\n")
}
