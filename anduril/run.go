package anduril

import "github.com/cicovic-andrija/anduril/service"

func Run() {
	env, err := service.ReadEnvironment()
	if err != nil {
		panic(err)
	}

	err = env.Initialize()
	if err != nil {
		panic(err)
	}

	config, err := ReadConfig(env.ConfigPath())
	if err != nil {
		panic(err)
	}

	server, err := NewWebServer(env, config)
	if err != nil {
		panic(err)
	}

	server.ListenAndServe()
}
