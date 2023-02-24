package anduril

func Run() {
	env, err := readEnv()
	if err != nil {
		panic(err)
	}

	err = initEnv(env)
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
