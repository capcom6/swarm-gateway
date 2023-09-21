package common

type Auth struct {
	Type string
	Data string
}

type Service struct {
	ID      string
	Version uint64
	Name    string
	Host    string
	Port    uint16
	Auth    Auth
}
