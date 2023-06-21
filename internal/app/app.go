package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/capcom6/swarm-gateway-tutorial/internal/discovery"
	"github.com/capcom6/swarm-gateway-tutorial/internal/repository"
	"github.com/docker/docker/client"
)

func Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup

	servicesRepo := repository.NewServicesRepository()
	if err := startDiscovery(ctx, &wg, servicesRepo); err != nil {
		return err
	}
	if err := startProxy(ctx, &wg, servicesRepo); err != nil {
		return err
	}

	wg.Wait()

	log.Println("Done")

	return nil
}

func startDiscovery(ctx context.Context, wg *sync.WaitGroup, servicesRepo *repository.ServicesRepository) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		discoverySvc := discovery.NewSwarmDiscovery(cli)
		timer := time.NewTicker(5 * time.Second)
		defer func() {
			timer.Stop()
			cli.Close()
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				log.Println("Discovery Done")
				return
			case <-timer.C:
				timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
				services, err := discoverySvc.ListServices(timeoutCtx)
				if err != nil {
					log.Println(err)
				}
				cancel()

				servicesRepo.ReplaceServices(services)
			}
		}
	}()

	return nil
}

func startProxy(ctx context.Context, wg *sync.WaitGroup, servicesRepo *repository.ServicesRepository) error {
	wg.Add(1)
	go func() {
		timer := time.NewTicker(5 * time.Second)
		defer func() {
			timer.Stop()
			wg.Done()
		}()

		log.Println("Proxy Started")
		for {
			select {
			case <-ctx.Done():
				log.Println("Proxy Done")
				return
			case <-timer.C:
				service, err := servicesRepo.GetServiceByHost("test.example.com")
				if err != nil {
					log.Println(err)
					continue
				}

				log.Printf("%s - %s:%d", service.Name, service.Host, service.Port)
			}
		}
	}()

	return nil
}
