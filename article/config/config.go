package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config contains the configuration of the server and database
type Config struct {
	Server Server `json:"server"`
	DB     DB     `json:"db"`
}

// Server configuration
type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// DB configuration
type DB struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

//  FromFile return a configuration from a given file
func FromFile(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg *Config

	// Unmarshal data
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

