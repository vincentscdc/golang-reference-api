package configuration

import "time"

type Config struct {
	Application struct {
		Version   string `yaml:"version"`
		Port      int    `yaml:"port"`
		PrettyLog bool   `yaml:"prettylog"`
		URL       struct {
			Host    string   `yaml:"host"`
			Schemes []string `yaml:"schemes"`
		} `yaml:"url"`
		Timeouts struct {
			ReadTimeout       time.Duration `yaml:"readTimeout"`
			ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
			WriteTimeout      time.Duration `yaml:"writeTimeout"`
			IdleTimeout       time.Duration `yaml:"idleTimeout"`
		}
	} `yaml:"application"`
	Grpc struct {
		Port int `yaml:"port"`
	} `yaml:"grpc"`
	Observability struct {
		Collector struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"collector"`
	} `yaml:"observability"`
	DB Database `yaml:"db"`
}

type Database struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	User         string        `yaml:"user"`
	Password     string        `yaml:"password"`
	Database     string        `yaml:"database"`
	MaxConns     int32         `yaml:"maxConns"`
	MaxIdleConns int32         `yaml:"maxIdleConns"`
	MaxLifeTime  time.Duration `yaml:"maxLifeTime"`
}
