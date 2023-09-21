project_name = swarm-gateway
image_name = capcom6/swarm-gateway
image_tag = dev

init:
	go mod download \
		&& go install github.com/cosmtrek/air@latest

run:
	go run cmd/$(project_name)/main.go

air:
	air

test:
	go test ./...

docker-build-amd64:
	docker buildx build --platform linux/amd64 -t $(image_name):$(image_tag) --build-arg PROJECT_NAME=$(project_name) -f build/package/Dockerfile .

docker-build:
	docker build -t $(image_name):$(image_tag) --build-arg PROJECT_NAME=$(project_name) -f build/package/Dockerfile .

.PHONY: init run air test docker-build docker-build-amd64
