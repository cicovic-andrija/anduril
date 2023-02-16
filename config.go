package main

import "github.com/cicovic-andrija/https"

type Config struct {
	HTTPS https.Config `json:"https"`
}
