package anduril

import (
	"fmt"

	"github.com/cicovic-andrija/anduril/service"
	"github.com/cicovic-andrija/go-util"
)

type TraceTag string

const (
	MarkdownProcessorTag   TraceTag = "MarkdownProcessor"
	RepositoryProcessorTag TraceTag = "RepositoryProcessor"
	ExecutorTag            TraceTag = "Executor"
	CleanupTaskTag         TraceTag = "CleanupTask"
)

func (s *WebServer) log(format string, v ...interface{}) {
	s.logger.Output(util.SevInfo, 2, format, v...)
}

func (s *WebServer) trace(tag TraceTag, format string, v ...interface{}) {
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

func (s *WebServer) generateTraceCallback(tag TraceTag) service.TraceCallback {
	return func(format string, v ...interface{}) {
		s.logger.Output(util.SevInfo, 2, "["+string(tag)+"]: "+format, v...)
	}
}
