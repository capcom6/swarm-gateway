package config

type Storage struct {
	Driver   string            `yaml:"driver"`
	Host     string            `yaml:"host"`
	Port     uint16            `yaml:"port"`
	User     string            `yaml:"user"`
	Password string            `yaml:"password"`
	Options  map[string]string `yaml:"options"`
}
