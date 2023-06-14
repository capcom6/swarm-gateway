package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/capcom6/swarm-gateway-tutorial/internal/discovery"
	"github.com/docker/docker/client"
)

func Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup

	if err := startDiscovery(ctx, &wg); err != nil {
		return err
	}

	wg.Wait()

	log.Println("Done")

	return nil
}

func startDiscovery(ctx context.Context, wg *sync.WaitGroup) error {
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

				for _, service := range services {
					log.Printf("%s - %s:%d", service.Name, service.Host, service.Port)
				}
			}
		}
	}()

	return nil
}
