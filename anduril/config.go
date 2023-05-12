package anduril

import (
	"encoding/json"

	"github.com/cicovic-andrija/anduril/repository"
	"github.com/cicovic-andrija/go-util"
	"github.com/cicovic-andrija/https"
)

type Config struct {
	HTTPS      https.Config      `json:"https"`
	Repository repository.Config `json:"repository"`
}

func ReadConfig(path string) (*Config, error) {
	config := &Config{}
	if configFile, err := util.OpenFile(path); err != nil {
		return nil, err
	} else {
		defer configFile.Close()
		if err = json.NewDecoder(configFile).Decode(config); err != nil {
			return nil, err
		}
	}
	return config, nil
}
