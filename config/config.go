package config

import "github.com/jinzhu/configor"

type Config struct {
	RoomProxy struct {
		Servers []struct {
			Name string
			Host string `default:"0.0.0.0"`
			Port int
		}
	}
}

func Load(path ...string) (*Config, error) {
	c := &Config{}
	cc := &configor.Config{Debug: false, Verbose: false}
	err := configor.New(cc).Load(c, path...)
	if err != nil {
		return nil, err
	}
	return c, nil
}
