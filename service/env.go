package service

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/cicovic-andrija/libgo/fs"
)

// External programs.
const (
	MarkdownHTMLConverter = "pandoc"
)

// Command-line options and their values.
const (
	ConfigOption    = "config"
	PlaintextOption = "plaintext"
)

// Version and build.
// TODO: Thid could be set at link-time but -X linker option wasn't doing the trick. Investigate.
var (
	Version = "v1.0.3-7841277"
	Build   = "d3c24a79-c533-4e2d-974a-b4aab92198a6"
)

type Environment struct {
	initd           bool
	pid             int
	wd              string
	configPath      string
	encryptedConfig bool
}

func ReadEnvironment() (*Environment, error) {
	env := &Environment{
		initd: false,
	}

	if runtime.GOOS != "linux" {
		return nil, errors.New("unsupported OS")
	}

	env.pid = os.Getpid()
	if exe, err := os.Executable(); err == nil {
		env.wd = filepath.Dir(exe)
	} else {
		return nil, err
	}

	if err := env.parseCommandLine(); err != nil {
		return nil, fmt.Errorf("invalid argument: %v", err)
	}

	if _, err := exec.LookPath(MarkdownHTMLConverter); err != nil {
		return nil, fmt.Errorf("dependency not found on the system: %s", MarkdownHTMLConverter)
	}

	if exists, _ := fs.DirectoryExists(env.DataDirectoryPath()); !exists {
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
		if err := fs.MkdirIfNotExists(directory); err != nil {
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

func (env *Environment) PrimaryLogPath() string {
	return filepath.Join(env.LogsDirectoryPath(), "anduril.log")
}

func (env *Environment) TemplatePath(templateName string) string {
	return filepath.Join(env.DataDirectoryPath(), "templates", templateName)
}

func (env *Environment) AssetsDataDirectory() string {
	return filepath.Join(env.DataDirectoryPath(), "assets")
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

func (env *Environment) parseCommandLine() error {
	var (
		plaintext bool
	)

	flag.StringVar(
		&env.configPath,
		ConfigOption,
		"",
		"config file full path",
	)

	flag.BoolVar(
		&plaintext,
		PlaintextOption,
		false,
		"plaintext config indicator",
	)

	flag.Parse()

	env.encryptedConfig = !plaintext

	if env.configPath == "" {
		if env.encryptedConfig {
			env.configPath = filepath.Join(env.DataDirectoryPath(), "encrypted-config.txt")
		} else {
			env.configPath = filepath.Join(env.DataDirectoryPath(), "anduril-config.json")
		}
	}

	return nil
}
