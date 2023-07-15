package config

var DefaultConfig = Config{
	Acme: Acme{
		DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory",
		Storage: Storage{
			Driver: "filesystem",
			Options: map[string]string{
				"path": "./certs",
			},
		},
	},
}
