package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case tm := <-timer.C:
				fmt.Printf("%s\n", tm)
				listServices(cli)
			}
		}
	}()

	<-ctx.Done()

	log.Println("Done")
}

func listServices(cli *client.Client) {
	list, err := cli.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		panic(err)
	}

	for _, service := range list {
		fmt.Println(service.Spec.Name)
		fmt.Println(service.Spec.Labels)
	}
}
