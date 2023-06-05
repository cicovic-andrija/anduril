package anduril

import (
	"fmt"

	"github.com/cicovic-andrija/anduril/service"
	"github.com/cicovic-andrija/libgo/logging"
)

type TraceTag string

const (
	MarkdownProcessorTag TraceTag = "MarkdownProcessor"
	RepositoryTag        TraceTag = "Repository"
	ExecutorTag          TraceTag = "Executor"
	CleanupTag           TraceTag = "Cleanup"
)

func (s *WebServer) log(format string, v ...interface{}) {
	s.logger.Output(logging.SevInfo, 2, format, v...)
}

func (s *WebServer) trace(tag TraceTag, format string, v ...interface{}) {
	s.logger.Output(logging.SevInfo, 2, "["+string(tag)+"]: "+format, v...)
}

func (s *WebServer) warn(format string, v ...interface{}) {
	s.logger.Output(logging.SevWarn, 2, format, v...)
}

func (s *WebServer) error(format string, v ...interface{}) error {
	err := fmt.Errorf(format, v...)
	s.logger.Output(logging.SevError, 2, format, v...)
	return err
}

func (s *WebServer) generateTraceCallback(tag TraceTag) service.TraceCallback {
	return func(format string, v ...interface{}) {
		s.logger.Output(logging.SevInfo, 2, "["+string(tag)+"]: "+format, v...)
	}
}
