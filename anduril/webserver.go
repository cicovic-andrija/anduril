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
	env            *Environment
	httpsServer    *https.HTTPSServer
	repository     *RepositoryProcessor
	latestRevision *Revision
	executor       *Executor
	logger         *util.FileLog
	startedAt      time.Time
}

func NewWebServer(env *Environment, config *Config) (server *WebServer, err error) {
	if env == nil || !env.initd {
		return nil, errors.New("environment not initialized")
	}

	logger, err := util.NewFileLog(env.PrimaryLogPath())
	if err != nil {
		return nil, fmt.Errorf("failed to create primary log file: %v", err)
	}

	config.HTTPS.LogsDirectory = env.LogsDirectoryPath()
	config.HTTPS.FileServer.Directory = filepath.Join(env.DataDirectoryPath(), staticSubdir)
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

	webServer.latestRevision = nil

	webServer.executor = &Executor{
		trace: webServer.generateTraceCallback(ExecutorTag),
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

	s.testGit()

	s.listenAndServeInternal()
}

func (s *WebServer) testGit() {
	err := s.repository.OpenOrClone(filepath.Join(s.env.WorkDirectoryPath(), repositorySubdir))
	if err != nil {
		s.error("%v", err)
		return
	}

	revision := &Revision{
		Articles:      make(map[string]*Article),
		Tags:          make(map[string][]*Article),
		ContainerPath: s.repository.ContentRoot(),
		Hash:          s.repository.LatestCommitShortHash(),
	}

	s.processRevision(revision)
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
