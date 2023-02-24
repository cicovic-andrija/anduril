package anduril

import (
	"fmt"
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
	if exists, err := util.DirectoryExists(env.logsDirPath()); err != nil {
		return fmt.Errorf("failed to stat logs directory: %v", err)
	} else if !exists {
		if err := os.Mkdir(env.logsDirPath(), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create logs directory: %v", err)
		}
	}

	env.initd = true
	return nil
}

func (env *envinfo) configPath() string {
	return filepath.Join(env.wd, "anduril-config.json")
}

func (env *envinfo) logsDirPath() string {
	return filepath.Join(env.wd, "logs")
}

func (env *envinfo) primaryLogPath() string {
	return filepath.Join(env.logsDirPath(), "anduril.log")
}
