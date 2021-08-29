service=tank

run:
	source .ENV && go run -race cmd/${service}/main.go

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go

docker:
	docker build -t ${service} . && docker run ${service}
