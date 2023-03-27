package anduril

import (
	"fmt"
	"io"
	"os/exec"
)

type Executor struct {
	trace TraceCallback
}

func (e *Executor) ConvertMarkdownToHTML(inputFilePath string, outputFilePath string) error {
	c := exec.Command(MarkdownHTMLConverter, "--from", "markdown", "--to", "html5", "--output", outputFilePath, inputFilePath)
	c.Stdout = io.Discard
	c.Stderr = io.Discard
	if err := c.Run(); err != nil {
		return fmt.Errorf("%s: %v", MarkdownHTMLConverter, err)
	}
	e.trace("%s: %s => %s", MarkdownHTMLConverter, inputFilePath, outputFilePath)
	return nil
}
