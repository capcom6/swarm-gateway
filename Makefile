project_name = swarm-gateway
image_name = $(project_name)
image_tag = latest

init:
	go mod download \
		&& go install github.com/cosmtrek/air@latest

run:
	go run cmd/$(project_name)/main.go

air:
	air

test:
	go test ./...

docker-build:
	docker build -t $(image_name):$(image_tag) --build-arg PROJECT_NAME=$(project_name) -f build/package/Dockerfile .

.PHONY: init run air test docker-build
