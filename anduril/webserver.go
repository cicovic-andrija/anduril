package anduril

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/cicovic-andrija/anduril/repository"
	"github.com/cicovic-andrija/anduril/service"
	"github.com/cicovic-andrija/go-util"
	"github.com/cicovic-andrija/https"
)

type WebServer struct {
	env            *service.Environment
	httpsServer    *https.HTTPSServer
	repository     repository.Repository
	latestRevision *Revision
	revisionLock   *sync.RWMutex
	executor       *Executor
	taskWaitGroup  *sync.WaitGroup
	stopChannels   []chan struct{}
	logger         *util.FileLog
	startedAt      time.Time
}

func NewWebServer(env *service.Environment, config *Config) (server *WebServer, err error) {
	if env == nil || !env.IsInitialized() {
		return nil, errors.New("environment not initialized")
	}

	if config == nil {
		return nil, errors.New("invalid configuration object: null")
	}

	logger, err := util.NewFileLog(env.PrimaryLogPath())
	if err != nil {
		return nil, fmt.Errorf("failed to create primary log file: %v", err)
	}

	config.HTTPS.LogsDirectory = env.LogsDirectoryPath()
	config.HTTPS.FileServer.Directory = env.StaticDataDirectory()
	httpsServer, err := https.NewServer(&config.HTTPS)
	if err != nil {
		return nil, fmt.Errorf("failed to init HTTPS server: %v", err)
	}

	webServer := &WebServer{
		env:         env,
		httpsServer: httpsServer,
		logger:      logger,
	}

	gitRepo := &repository.GitRepository{
		Config: config.Repository,
	}
	if err = gitRepo.Validate(); err != nil {
		return nil, fmt.Errorf("invalid repository configuration: %v", err)
	} else {
		webServer.repository = gitRepo
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
	s.log("pid: %d", s.env.PID())
	s.log("working directory: %s", s.env.WDP())
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
			s.stopPeriodicTasks()
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
	const N = 2
	s.taskWaitGroup = &sync.WaitGroup{}
	s.stopChannels = make([]chan struct{}, N)
	for i := range s.stopChannels {
		s.stopChannels[i] = make(chan struct{})
	}

	s.taskWaitGroup.Add(N)
	go s.genericPeriodicTask(s.checkForNewRevision, 15*time.Minute, s.stopChannels[0], RepositoryTag)
	go s.genericPeriodicTask(s.cleanUpCompiledFiles, 24*time.Hour, s.stopChannels[1], CleanupTag)
}

func (s *WebServer) stopPeriodicTasks() {
	for _, task := range s.stopChannels {
		close(task)
	}
	s.taskWaitGroup.Wait()
}

func (s *WebServer) genericPeriodicTask(f func(service.TraceCallback, ...interface{}) error, period time.Duration, stop chan struct{}, tag TraceTag, v ...interface{}) {
	s.log("starting timer task %s", tag)
	trace := s.generateTraceCallback(tag)
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ticker.C:
			err := f(trace, v...)
			if err != nil {
				s.error("timer task %s failed with error: %v", tag, err)
			}
		case <-stop:
			s.taskWaitGroup.Done()
			s.log("timer task %s stopped", tag)
			return
		}
	}
}
