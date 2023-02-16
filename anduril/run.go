package anduril

import (
	"os"
	"path/filepath"

	"github.com/cicovic-andrija/https"
)

func Run() {
	var (
		wd string
	)

	if exe, err := os.Executable(); err == nil {
		wd = filepath.Dir(exe)
	} else {
		panic(err)
	}

	config, err := readConfig(filepath.Join(wd, "anduril-config.json"))
	if err != nil {
		panic(err)
	}

	httpsServer, err := https.NewServer(&config.HTTPS)
	if err != nil {
		panic(err)
	}

	errorChannel := make(chan error, 1)
	interruptChannel := make(chan os.Signal, 1)
	httpsServer.ListenAndServeAsync(errorChannel)
	for {
		select {
		case <-interruptChannel:
			if shutdownError := httpsServer.Shutdown(); err != nil {
				// TODO: Log shutdownError
				panic(shutdownError)
			}
			os.Exit(0)
		case serverError := <-errorChannel:
			// TODO: Log serverError
			panic(serverError)
		}
	}
}
