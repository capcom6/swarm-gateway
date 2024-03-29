package acme

import (
	"context"
	"crypto/tls"

	"github.com/capcom6/swarm-gateway/internal/config"
	"github.com/capcom6/swarm-gateway/internal/repository"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func NewConfig(services *repository.ServicesRepository, cache autocert.Cache, acmeCfg config.Acme) *tls.Config {
	// Certificate manager
	m := &autocert.Manager{
		Client: makeClient(acmeCfg),
		Prompt: autocert.AcceptTOS,
		// Replace with your domain
		HostPolicy: func(ctx context.Context, host string) error {
			_, err := services.GetServiceByHost(host)
			return err
		},
		// Folder to store the certificates
		Cache: cache,
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

func makeClient(cfg config.Acme) *acme.Client {
	return &acme.Client{
		DirectoryURL: cfg.DirectoryURL,
	}
}
