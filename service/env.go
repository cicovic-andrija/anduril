package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/cicovic-andrija/go-util"
)

// External programs.
const (
	MarkdownHTMLConverter = "pandoc"
)

type Environment struct {
	wd    string
	pid   int
	initd bool
}

func ReadEnvironment() (*Environment, error) {
	env := &Environment{
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

	if exists, _ := util.DirectoryExists(env.DataDirectoryPath()); !exists {
		return env, fmt.Errorf("directory not found: %s", env.DataDirectoryPath())
	}

	return env, nil
}

func (env *Environment) Initialize() error {
	for _, directory := range []string{
		env.WorkDirectoryPath(),
		env.LogsDirectoryPath(),
		env.RepositoryWorkingDirectory(),
		env.CompiledWorkDirectory(),
	} {
		if err := util.MkdirIfNotExists(directory); err != nil {
			return err
		}
	}

	env.initd = true
	return nil
}

func (env *Environment) IsInitialized() bool {
	return env.initd
}

func (env *Environment) PID() int {
	return env.pid
}

func (env *Environment) WDP() string {
	return env.wd
}

func (env *Environment) DataDirectoryPath() string {
	return filepath.Join(env.wd, "data")
}

func (env *Environment) WorkDirectoryPath() string {
	return filepath.Join(env.wd, "work")
}

func (env *Environment) LogsDirectoryPath() string {
	return filepath.Join(env.wd, "logs")
}

func (env *Environment) ConfigPath() string {
	return filepath.Join(env.DataDirectoryPath(), "anduril-config.json")
}

func (env *Environment) PrimaryLogPath() string {
	return filepath.Join(env.LogsDirectoryPath(), "anduril.log")
}

func (env *Environment) TemplatePath(templateName string) string {
	return filepath.Join(env.DataDirectoryPath(), "templates", templateName)
}

func (env *Environment) StaticDataDirectory() string {
	return filepath.Join(env.DataDirectoryPath(), "static")
}

func (env *Environment) RepositoryWorkingDirectory() string {
	return filepath.Join(env.WorkDirectoryPath(), "repository")
}

func (env *Environment) CompiledWorkDirectory() string {
	return filepath.Join(env.WorkDirectoryPath(), "compiled")
}

func (env *Environment) CompiledTemplatePath(templateName string) string {
	return filepath.Join(env.CompiledWorkDirectory(), templateName)
}
