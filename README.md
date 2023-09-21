# Swarm Gateway

Swarm Gateway is a reverse-proxy application designed specifically for Docker Swarm clusters. It utilizes Docker Swarm service labels to configure proxying.

The project started as tutorial for co-workers and now is in the MVP stage. Not for production use.

## Features

* Automatic proxy configuration based on Docker Swarm service labels
* Support for dynamically scaling services in Docker Swarm
* Acme challenge support
* Authentication support

## Getting Started

To get started with Swarm Gateway, follow these steps:

1. Install Docker Swarm on your cluster.
2. Deploy your services to the Docker Swarm cluster, making sure to include the necessary labels for proxy configuration.
3. Run the Swarm Gateway app on a manager node in the Swarm cluster.

```
docker service create \
  --name swarm-gateway \
  --replicas 1 \
  -p 80:3000 \
  -p 443:3443 \
  --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
  capcom6/swarm-gateway:latest
```

This command creates a Docker service named swarm-gateway with 1 replica, exposing ports 80 and 443 for HTTP and HTTPS traffic respectively. It also mounts the Docker socket to enable communication with the Docker API.

Swarm Gateway will automatically detect the services and their labels, and configure the reverse proxy accordingly.

## Configuration

Swarm Gateway uses the following Docker Swarm service labels for configuration:

* gateway.enabled: Set this label to true to enable proxying for the service.
* gateway.server.port: The port on which the service is running.
* gateway.server.host: External host name of the service.
* gateway.auth.type: The type of authentication to use (e.g., `basic`).
* gateway.auth.data: The authentication data required for the selected authentication type.

### Example

```
docker service create \
  --name my-service \
  --label gateway.enabled=true \
  --label gateway.server.port=8080 \
  --label gateway.server.host=service.example.com \
  --label gateway.auth.type=basic \
  --label gateway.auth.data=demo:$2y$05$il5tgnyhkmzOBVcpLI2greDYn/m2.ja75dEqVyMRX/P/xQfIhVvZC \
  my-service-image
```

## Roadmap

* Support for custom routing rules based on additional service labels
* Metrics and monitoring for the proxy traffic
* Support for WebSocket communication
* Fine-grained access control and rate limiting

## Contributing
Contributions are welcome! If you have any suggestions, bug reports, or feature requests, please open an issue or submit a pull request.

## License

This project is licensed under the Apache License 2.0. See the LICENSE file for details.

## Acknowledgments

This project was inspired by Traefik and aims to provide a simplified reverse-proxy solution for Docker Swarm clusters.