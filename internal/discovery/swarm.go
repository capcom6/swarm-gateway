package discovery

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/capcom6/swarm-gateway/internal/common"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

const (
	LabelGatewayEnabled    = "gateway.enabled"
	LabelGatewayServerPort = "gateway.server.port"
	LabelGatewayServerHost = "gateway.server.host"
	LabelGatewayAuthType   = "gateway.auth.type"
	LabelGatewayAuthData   = "gateway.auth.data"

	NetworkName = "proxy"
)

type SwarmDiscovery struct {
	Client *client.Client
}

func NewSwarmDiscovery(client *client.Client) *SwarmDiscovery {
	return &SwarmDiscovery{
		Client: client,
	}
}

func (sd *SwarmDiscovery) ListServices(ctx context.Context) ([]common.Service, error) {
	networkId, err := sd.getNetworkIdByName(ctx, NetworkName)
	if err != nil {
		return nil, err
	}

	list, err := sd.Client.ServiceList(ctx, types.ServiceListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", LabelGatewayEnabled+"=true"),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("can't list services: %w", err)
	}

	services := make([]common.Service, 0, len(list))
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

		if !sd.isServiceConnectedToNetwork(networkId, service) {
			continue
		}

		services = append(services, common.Service{
			ID:      service.ID,
			Version: service.Version.Index,
			Name:    service.Spec.Name,
			Host:    service.Spec.Labels[LabelGatewayServerHost],
			Port:    uint16(port),
			Auth: common.Auth{
				Type: service.Spec.Labels[LabelGatewayAuthType],
				Data: service.Spec.Labels[LabelGatewayAuthData],
			},
		})
	}

	return services, nil
}

func (sd *SwarmDiscovery) getNetworkIdByName(ctx context.Context, name string) (string, error) {
	networks, err := sd.Client.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(
			filters.Arg("name", name),
		),
	})
	if err != nil {
		return "", fmt.Errorf("can't list networks: %w", err)
	}

	if len(networks) == 0 {
		return "", fmt.Errorf("network %s not found", name)
	}

	return networks[0].ID, nil
}

func (sd *SwarmDiscovery) isServiceConnectedToNetwork(id string, service swarm.Service) bool {
	for _, network := range service.Spec.TaskTemplate.Networks {
		if network.Target == id {
			return true
		}
	}
	return false
}
