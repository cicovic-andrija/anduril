package anduril

import (
	"os"
	"path/filepath"

	"github.com/cicovic-andrija/go-util"
)

type envinfo struct {
	wd    string
	pid   int
	initd bool
}

func readEnv() (*envinfo, error) {
	env := &envinfo{
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

func initEnv(env *envinfo) error {
	if err := util.MkdirIfNotExists(env.workDirPath()); err != nil {
		return err
	}

	if err := util.MkdirIfNotExists(env.logsDirPath()); err != nil {
		return err
	}

	env.initd = true
	return nil
}

func (env *envinfo) configPath() string {
	return filepath.Join(env.wd, "anduril-config.json")
}

func (env *envinfo) workDirPath() string {
	return filepath.Join(env.wd, "work")
}

func (env *envinfo) logsDirPath() string {
	return filepath.Join(env.wd, "logs")
}

func (env *envinfo) primaryLogPath() string {
	return filepath.Join(env.logsDirPath(), "anduril.log")
}
