package config

type Storage struct {
	Driver  string            `yaml:"driver"`
	Options map[string]string `yaml:"options"`
}
