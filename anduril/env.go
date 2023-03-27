package anduril

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/cicovic-andrija/go-util"
)

// Data subdirectories.
const (
	staticSubdir    = "static"
	templatesSubdir = "templates"
)

// Work subdirectories.
const (
	repositorySubdir = "repository"
	compiledSubdir   = "compiled"
)

// External programs.
const (
	MarkdownHTMLConverter = "pandoc"
)

type environment struct {
	wd    string
	pid   int
	initd bool
}

func readEnvironment() (*environment, error) {
	env := &environment{
		initd: false,
	}

	if runtime.GOOS != "linux" {
		return nil, errors.New("unsupported OS")
	}

	if _, err := exec.LookPath(MarkdownHTMLConverter); err != nil {
		return nil, fmt.Errorf("dependency not found on the system: %s", MarkdownHTMLConverter)
	}

	env.pid = os.Getpid()
	if exe, err := os.Executable(); err == nil {
		env.wd = filepath.Dir(exe)
	} else {
		return nil, err
	}

	if exists, _ := util.DirectoryExists(env.dataDirectoryPath()); !exists {
		return env, fmt.Errorf("failed to stat dir: %s", env.dataDirectoryPath())
	}

	return env, nil
}

func prepareEnvironment(env *environment) error {
	for _, directory := range []string{
		env.workDirectoryPath(),
		filepath.Join(env.workDirectoryPath(), repositorySubdir),
		filepath.Join(env.workDirectoryPath(), compiledSubdir),
		filepath.Join(env.workDirectoryPath(), staticSubdir),
		env.logsDirectoryPath(),
	} {
		if err := util.MkdirIfNotExists(directory); err != nil {
			return err
		}
	}

	env.initd = true
	return nil
}

func (env *environment) dataDirectoryPath() string {
	return filepath.Join(env.wd, "data")
}

func (env *environment) workDirectoryPath() string {
	return filepath.Join(env.wd, "work")
}

func (env *environment) logsDirectoryPath() string {
	return filepath.Join(env.wd, "logs")
}

func (env *environment) configPath() string {
	return filepath.Join(env.wd, "anduril-config.json")
}

func (env *environment) primaryLogPath() string {
	return filepath.Join(env.logsDirectoryPath(), "anduril.log")
}
