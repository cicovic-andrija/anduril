package anduril

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/cicovic-andrija/go-util"
	"github.com/cicovic-andrija/https"
)

type WebServer struct {
	env            *Environment
	httpsServer    *https.HTTPSServer
	repository     *RepositoryProcessor
	latestRevision *Revision
	revisionLock   *sync.RWMutex
	executor       *Executor
	taskWaitGroup  *sync.WaitGroup
	stopChannels   []chan struct{}
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
		env:           env,
		httpsServer:   httpsServer,
		taskWaitGroup: &sync.WaitGroup{},
		stopChannels:  make([]chan struct{}, 0),
		logger:        logger,
	}

	// Explicitly set repo to nil: struct not initialized.
	webServer.repository = &RepositoryProcessor{
		RepositoryConfig: config.Repository,
		repo:             nil,
	}

	// Explicitly set to nil: not initialized.
	webServer.latestRevision = nil
	webServer.revisionLock = &sync.RWMutex{}

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
	s.startPeriodicTasks()
	s.listenAndServeInternal()
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

func (s *WebServer) startPeriodicTasks() {
	s.taskWaitGroup.Add(1)
	stop := make(chan struct{})
	go s.genericPeriodicTask(s.checkForNewRevision, 10*time.Second, stop, RepositoryProcessorTag)
}

func (s *WebServer) genericPeriodicTask(f func(TraceCallback, ...interface{}) error, period time.Duration, stop chan struct{}, tag TraceTag, v ...interface{}) {
	s.log("starting timer task %s", tag)
	trace := s.generateTraceCallback(tag)
	ticker := time.NewTicker(period)
	for {
		select {
		case <-stop:
			s.taskWaitGroup.Done()
			s.log("timer task %s stopped", tag)
		case <-ticker.C:
			err := f(trace, v...)
			if err != nil {
				s.error("timer task %s failed with error: %v", tag, err)
			}
		}
	}
}
