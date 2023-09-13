package config

type Acme struct {
	DirectoryURL string  `yaml:"directoryUrl"`
	Email        string  `yaml:"email"`
	Storage      Storage `yaml:"storage"`
}
