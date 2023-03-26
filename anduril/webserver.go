package anduril

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/cicovic-andrija/go-util"
	"github.com/cicovic-andrija/https"
)

type WebServer struct {
	env         *environment
	httpsServer *https.HTTPSServer
	repository  *RepositoryProcessor
	logger      *util.FileLog
	startedAt   time.Time
}

func NewWebServer(env *environment, config *Config) (server *WebServer, err error) {
	if env == nil || !env.initd {
		return nil, errors.New("environment not initialized")
	}

	logger, err := util.NewFileLog(env.primaryLogPath())
	if err != nil {
		return nil, fmt.Errorf("failed to create primary log file: %v", err)
	}

	config.HTTPS.LogsDirectory = env.logsDirectoryPath()
	config.HTTPS.FileServer.Directory = filepath.Join(env.workDirectoryPath(), staticSubdir)
	httpsServer, err := https.NewServer(&config.HTTPS)
	if err != nil {
		return nil, fmt.Errorf("failed to init HTTPS server: %v", err)
	}

	webServer := &WebServer{
		env:         env,
		httpsServer: httpsServer,
		logger:      logger,
	}

	webServer.repository = &RepositoryProcessor{
		RepositoryConfig: config.Repository,
		trace:            webServer.generateTraceCallback(RepositoryProcessorTag),
		repo:             nil,
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

func (s *WebServer) prepareEnvironment() {

}

func (s *WebServer) listenAndServeInternal() {
	// First, register handlers.
	s.registerHandlers()

	// Start accepting HTTPS connections.
	httpsErrorChannel := make(chan error, 1)
	s.httpsServer.ListenAndServeAsync(httpsErrorChannel)

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt)
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
