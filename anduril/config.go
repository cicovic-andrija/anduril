package anduril

import (
	"encoding/json"

	"github.com/cicovic-andrija/go-util"
	"github.com/cicovic-andrija/https"
)

type Config struct {
	HTTPS      https.Config     `json:"https"`
	Repository RepositoryConfig `json:"repository"`
}

type RepositoryConfig struct {
	URL                 string `json:"url"`
	Remote              string `json:"remote"`
	Branch              string `json:"branch"`
	RelativeContentPath string `json:"relative_content_path"`
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
