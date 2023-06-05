package anduril

import (
	"github.com/cicovic-andrija/anduril/service"
)

func Run() {
	env, err := service.ReadEnvironment()
	if err != nil {
		panic(err)
	}

	err = env.Initialize()
	if err != nil {
		panic(err)
	}

	config := &Config{}
	err = env.UnmarshalConfig(config)
	if err != nil {
		panic(err)
	}

	server, err := NewWebServer(env, config)
	if err != nil {
		panic(err)
	}

	server.ListenAndServe()
}
