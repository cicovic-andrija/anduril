package anduril

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/cicovic-andrija/go-util"
	"github.com/cicovic-andrija/https"
)

type WebServer struct {
	env         *envinfo
	httpsServer *https.HTTPSServer
	taskManager *util.TaskManager
	logger      *util.FileLog
	startedAt   time.Time
}

func NewWebServer(env *envinfo, config *Config) (server *WebServer, err error) {
	if env == nil || !env.initd {
		return nil, errors.New("environment not initialized")
	}

	logger, err := util.NewFileLog(env.primaryLogPath())
	if err != nil {
		return nil, fmt.Errorf("failed to create primary log file: %v", err)
	}

	config.HTTPS.LogsDirectory = env.logsDirPath()
	httpsServer, err := https.NewServer(&config.HTTPS)
	if err != nil {
		return nil, fmt.Errorf("failed to init HTTPS server: %v", err)
	}

	webServer := &WebServer{
		env:         env,
		httpsServer: httpsServer,
		taskManager: util.NewTaskManager(&logger.Logger),
		logger:      logger,
	}

	return webServer, nil
}

func (s *WebServer) ListenAndServe() {
	s.startedAt = time.Now().UTC()
	s.log("pid: %d", s.env.pid)
	s.log("working directory: %s", s.env.wd)
	s.log("primary log location: %s", s.logger.LogPath())
	s.log("HTTPS server log location: %s", s.httpsServer.GetLogPath())
	s.log("HTTPS requests log location: %s", s.httpsServer.GetRequestsLogPath())
	s.listenAndServeInternal()
}

func (s *WebServer) listenAndServeInternal() {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt)

	// Start accepting HTTPS connections.
	httpsErrorChannel := make(chan error, 1)
	s.httpsServer.ListenAndServeAsync(httpsErrorChannel)

	for {
		select {
		case <-interruptChannel:
			s.log("interrupt signal received, shutting down...")
			if httpsShutdownError := s.httpsServer.Shutdown(); httpsShutdownError != nil {
				panic(s.error("server shutdown: error encountered: %v", httpsShutdownError))
			}
			s.log("server was successfully shut down")
			os.Exit(0)
		case httpsServerError := <-httpsErrorChannel:
			panic(s.error("HTTPS server stopped unexpectedly:  %v", httpsServerError))
		}
	}
}

func (s *WebServer) log(format string, v ...interface{}) {
	s.logger.Output(util.SevInfo, 2, format, v...)
}

func (s *WebServer) logwarn(format string, v ...interface{}) {
	s.logger.Output(util.SevWarn, 2, format, v...)
}

func (s *WebServer) error(format string, v ...interface{}) error {
	err := fmt.Errorf(format, v...)
	s.logger.Output(util.SevError, 2, format, v...)
	return err
}
