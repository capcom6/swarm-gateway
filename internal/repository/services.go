package repository

import (
	"errors"
	"sync"

	"github.com/capcom6/swarm-gateway-tutorial/internal/common"
)

var ErrSeviceNotFound = errors.New("service not found")

type ServicesRepository struct {
	services map[string]common.Service
	mux      sync.RWMutex
}

func NewServicesRepository() *ServicesRepository {
	return &ServicesRepository{
		services: make(map[string]common.Service),
	}
}

func (sr *ServicesRepository) ReplaceServices(services []common.Service) {
	sr.mux.Lock()
	defer sr.mux.Unlock()

	sr.services = make(map[string]common.Service, len(services))
	for _, service := range services {
		sr.services[service.Host] = service
	}
}

func (sr *ServicesRepository) GetServiceByHost(host string) (common.Service, error) {
	sr.mux.RLock()
	defer sr.mux.RUnlock()

	if service, ok := sr.services[host]; ok {
		return service, nil
	}
	return common.Service{}, ErrSeviceNotFound
}
