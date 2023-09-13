package config

import "time"

var DefaultConfig = Config{
	Acme: Acme{DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory", Storage: Storage{Driver: "filesystem", Options: map[string]string{"path": "./certs"}}},
	Proxy: Proxy{
		Timeout: 120 * time.Second,
	},
}
