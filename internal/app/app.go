package app

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/capcom6/swarm-gateway-tutorial/internal/common"
	"github.com/capcom6/swarm-gateway-tutorial/internal/config"
	"github.com/capcom6/swarm-gateway-tutorial/internal/discovery"
	"github.com/capcom6/swarm-gateway-tutorial/internal/proxy"
	"github.com/capcom6/swarm-gateway-tutorial/internal/proxy/acme"
	"github.com/capcom6/swarm-gateway-tutorial/internal/proxy/acme/cache"
	"github.com/capcom6/swarm-gateway-tutorial/internal/proxy/auth"
	"github.com/capcom6/swarm-gateway-tutorial/internal/proxy/resolver"
	"github.com/capcom6/swarm-gateway-tutorial/internal/repository"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/valyala/fasthttp"
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
	config := config.Get()
	app := fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"10.0.0.0/8"},
		ProxyHeader:             "X-Forwarded-For",
	})

	app.Use(logger.New(logger.Config{
		Format:     `${ip} - [${time}] "${method} ${path} HTTP/1.1" ${host} ${status} ${bytesSent} ${latency}` + "\n",
		TimeFormat: "2006/01/02 15:04:05",
	}))
	app.Use(recover.New())
	app.Use(resolver.New(servicesRepo))
	app.Use(auth.New())

	app.Use(func(c *fiber.Ctx) error {
		service := c.Locals("service").(common.Service)

		query := string(c.Context().URI().QueryString())
		url := fmt.Sprintf("http://%s:%d%s", service.Name, service.Port, c.Path())
		if len(query) > 0 {
			url += "?" + query
		}

		if err := proxy.DoTimeout(c, url, service.Host, config.Proxy.Timeout); err != nil {
			log.Printf("proxy error: %s", err)
			if errors.Is(err, fasthttp.ErrTimeout) {
				return fiber.ErrGatewayTimeout
			}
			return fiber.ErrBadGateway
		}
		// Remove Server header from response
		c.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	certsCache, err := cache.New(config.Acme.Storage)
	if err != nil {
		return fmt.Errorf("can't create acme cache: %w", err)
	}

	tlsListener, err := tls.Listen("tcp", ":3443", acme.NewConfig(servicesRepo, certsCache, config.Acme))
	if err != nil {
		return fmt.Errorf("can't listen: %w", err)
	}

	wg.Add(1)
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Printf("can't listen: %s", err)
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err := app.Listener(tlsListener); err != nil {
			log.Printf("can't listen: %s", err)
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		<-ctx.Done()

		app.Shutdown()

		wg.Done()
	}()

	return nil
}
