package cache

import (
	"errors"

	"github.com/capcom6/swarm-gateway/internal/config"
	"github.com/capcom6/swarm-gateway/internal/proxy/acme/cache/s3"
	"golang.org/x/crypto/acme/autocert"
)

func New(cfg config.Storage) (autocert.Cache, error) {
	if cfg.Driver == "filesystem" {
		return autocert.DirCache(cfg.Options["path"]), nil
	}
	if cfg.Driver == "s3" {
		return s3.New(cfg.Options)
	}

	return nil, errors.New("driver not supported")
}
