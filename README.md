# qooked

A food recipe app

## Running locally

Run the following command:

go run cmd/api/main.go

## Running locally with Docker

Run the following commands:

docker build -t qooked-api-image:latest .

docker run --name qooked-api -d qooked-api-image:latest
