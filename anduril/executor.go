package anduril

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/cicovic-andrija/anduril/service"
)

type Executor struct {
	trace service.TraceCallback
}

func (e *Executor) ConvertMarkdownToHTML(inputFilePath string, outputFilePath string) error {
	c := exec.Command(service.MarkdownHTMLConverter, "--from", "markdown", "--to", "html5", "--output", outputFilePath, inputFilePath)
	c.Stdout = io.Discard
	c.Stderr = io.Discard
	if err := c.Run(); err != nil {
		return fmt.Errorf("%s: %v", service.MarkdownHTMLConverter, err)
	}
	e.trace("%s: %s => %s", service.MarkdownHTMLConverter, inputFilePath, outputFilePath)
	return nil
}
