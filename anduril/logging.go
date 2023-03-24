package anduril

import (
	"fmt"

	"github.com/cicovic-andrija/go-util"
)

type CheckpointTag string

const (
	MarkdownProcessorTag CheckpointTag = "markdown processor"
)

func (s *WebServer) log(format string, v ...interface{}) {
	s.logger.Output(util.SevInfo, 2, format, v...)
}

func (s *WebServer) checkpoint(tag CheckpointTag, format string, v ...interface{}) {
	s.logger.Output(util.SevInfo, 2, "["+string(tag)+"]: "+format, v...)
}

func (s *WebServer) warn(format string, v ...interface{}) {
	s.logger.Output(util.SevWarn, 2, format, v...)
}

func (s *WebServer) error(format string, v ...interface{}) error {
	err := fmt.Errorf(format, v...)
	s.logger.Output(util.SevError, 2, format, v...)
	return err
}
