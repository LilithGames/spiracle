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
		Enable  bool `default:"true"`
		Debug   bool `default:"false"`
		Workers int  `default:"1"`
		Session struct {
			MaxIdleDuration int `default:"30"`
			Expire          int `default:"30"`
		}
		Servers []Server
		MetricsAddr string `default:":2222"`
	}
	Controller struct {
		Reconciler struct {
			Enable bool `default:"true"`
		}
		Port           int    `default:"9443"`
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
