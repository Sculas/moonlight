package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"log"
	"os"
)

type config struct {
	Server struct {
		// Port indicates which port the server should listen on.
		Port int
		// Multicore indicates whether the server should use multiple cores.
		Multicore bool
		// ReuseAddr indicates whether to set up the SO_REUSEADDR socket option.
		ReuseAddr bool
		// ReusePort indicates whether to set up the SO_REUSEPORT socket option.
		ReusePort bool
	}
}

func defaultConfig() config {
	c := config{}
	c.Server.Port = 25565
	c.Server.Multicore = true
	return c
}

var (
	Config     *config
	configFile = "config.toml"
)

func Initialize() {
	c, err := readConfig()
	if err != nil {
		log.Fatalln(err)
		return
	}
	Config = &c
}

func readConfig() (config, error) {
	c := defaultConfig()
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := os.WriteFile(configFile, data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}
