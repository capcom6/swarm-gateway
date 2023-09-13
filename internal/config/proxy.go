package config

import "time"

type Proxy struct {
	Timeout time.Duration `yaml:"timeout"`
}
