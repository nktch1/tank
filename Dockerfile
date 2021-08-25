FROM debian:stable-slim as container
RUN apt-get update && apt-get install -yqq ca-certificates

FROM golang:1.16 as dep
WORKDIR /tmp/build
COPY go.mod go.sum ./
RUN go mod download

FROM dep as builder
ARG service
COPY . .
RUN cd $service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go

FROM container
ARG service
WORKDIR /app
COPY --from=builder /tmp/build/${service}/app .
CMD ["./app"]
