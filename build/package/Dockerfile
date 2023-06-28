# Stage 1 - Build the Golang Fiber app
FROM golang:1.19-alpine AS builder

ARG PROJECT_NAME

RUN apk add --no-cache git
RUN apk add --no-cache build-base

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/main cmd/${PROJECT_NAME}/main.go

# Stage 2 - Create the final Docker image
FROM alpine:3

COPY --from=builder /app/main /app/main

EXPOSE 3000

CMD ["/app/main"]