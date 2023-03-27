package anduril

func Run() {
	env, err := ReadEnvironment()
	if err != nil {
		panic(err)
	}

	err = PrepareEnvironment(env)
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
