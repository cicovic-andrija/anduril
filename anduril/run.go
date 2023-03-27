package anduril

func Run() {
	env, err := readEnvironment()
	if err != nil {
		panic(err)
	}

	err = prepareEnvironment(env)
	if err != nil {
		panic(err)
	}

	config, err := readConfig(env.configPath())
	if err != nil {
		panic(err)
	}

	server, err := NewWebServer(env, config)
	if err != nil {
		panic(err)
	}

	server.ListenAndServe()
}
