module github.com/cicovic-andrija/anduril

go 1.20

require (
	github.com/cicovic-andrija/go-util v0.0.0
	github.com/cicovic-andrija/https v0.0.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/cicovic-andrija/go-util => ../go-util
	github.com/cicovic-andrija/https => ../https
)
