package config

import "github.com/jinzhu/configor"

type Server struct {
	Name      string   `required:"true"`
	Host      string   `default:"0.0.0.0"`
	Port      int      `required:"true"`
	Externals []string `required:"true"`
}

type Config struct {
	RoomProxy struct {
		Debug   bool `default:"false"`
		Workers int  `default:"1"`
		Servers []Server
	}
	Controller struct {
		MetricsAddr    string `default:":8080"`
		ProbeAddr      string `default:":8081"`
		LeaderElection struct {
			Enable bool   `default:"true"`
			Id     string `default:"default-election-id"`
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
