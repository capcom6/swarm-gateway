package listener

import (
	"context"
	"crypto/tls"

	"github.com/capcom6/swarm-gateway-tutorial/internal/repository"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func NewConfig(services *repository.ServicesRepository) *tls.Config {
	// Certificate manager
	m := &autocert.Manager{
		Client: &acme.Client{
			DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory",
		},
		Prompt: autocert.AcceptTOS,
		// Replace with your domain
		HostPolicy: func(ctx context.Context, host string) error {
			_, err := services.GetServiceByHost(host)
			return err
		},
		// Folder to store the certificates
		Cache: autocert.DirCache("./certs"),
	}

	// TLS Config
	cfg := tls.Config{
		// Get Certificate from Let's Encrypt
		GetCertificate: m.GetCertificate,
		// By default NextProtos contains the "h2"
		// This has to be removed since Fasthttp does not support HTTP/2
		// Or it will cause a flood of PRI method logs
		// http://webconcepts.info/concepts/http-method/PRI
		NextProtos: []string{
			"http/1.1", "acme-tls/1",
		},
	}
	return &cfg
}
