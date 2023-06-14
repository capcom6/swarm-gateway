package discovery

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const (
	LabelGatewayEnabled    = "gateway.enabled"
	LabelGatewayServerPort = "gateway.server.port"
	LabelGatewayServerHost = "gateway.server.host"
)

type SwarmDiscovery struct {
	Client *client.Client
}

func NewSwarmDiscovery(client *client.Client) *SwarmDiscovery {
	return &SwarmDiscovery{
		Client: client,
	}
}

func (sd *SwarmDiscovery) ListServices(ctx context.Context) ([]Service, error) {
	list, err := sd.Client.ServiceList(ctx, types.ServiceListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", LabelGatewayEnabled+"=true"),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("can't list services: %w", err)
	}

	services := make([]Service, 0, len(list))
	for _, service := range list {
		host := service.Spec.Labels[LabelGatewayServerHost]
		if host == "" {
			log.Printf("no hostname for %s", service.Spec.Name)
			continue
		}

		port, err := strconv.Atoi(service.Spec.Labels[LabelGatewayServerPort])
		if err != nil {
			log.Printf("can't parse port for %s: %s", service.Spec.Name, err)
			continue
		}

		services = append(services, Service{
			ID:   service.ID,
			Name: service.Spec.Name,
			Host: service.Spec.Labels[LabelGatewayServerHost],
			Port: uint16(port),
		})
	}

	return services, nil
}
