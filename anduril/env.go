package anduril

import (
	"os"
	"path/filepath"

	"github.com/cicovic-andrija/go-util"
)

const (
	repositorySubdir = "repository"
	compiledSubdir   = "compiled"
	staticSubdir     = "static"
)

type environment struct {
	wd    string
	pid   int
	initd bool
}

func readEnv() (*environment, error) {
	env := &environment{
		initd: false,
	}

	if exe, err := os.Executable(); err == nil {
		env.wd = filepath.Dir(exe)
	} else {
		return nil, err
	}

	env.pid = os.Getpid()

	return env, nil
}

func initEnv(env *environment) error {
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

func (env *environment) configPath() string {
	return filepath.Join(env.wd, "anduril-config.json")
}

func (env *environment) workDirectoryPath() string {
	return filepath.Join(env.wd, "work")
}

func (env *environment) logsDirectoryPath() string {
	return filepath.Join(env.wd, "logs")
}

func (env *environment) primaryLogPath() string {
	return filepath.Join(env.logsDirectoryPath(), "anduril.log")
}
