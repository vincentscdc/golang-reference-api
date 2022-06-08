package configuration

import "time"

type Config struct {
	Application struct {
		Port      int
		PrettyLog bool
		URL       struct {
			Host    string
			Schemes []string
		}
		Timeouts struct {
			ReadTimeout       time.Duration
			ReadHeaderTimeout time.Duration
			WriteTimeout      time.Duration
			IdleTimeout       time.Duration
		}
	}
	Observability struct {
		Collector struct {
			Host string
			Port int
		}
	}
}
