tag=$(shell git describe --tags --always)
env=CGO_ENABLED=0 GOOS=linux GOARCH=amd64

build:
	echo "build $(tag)"
	@$(env) go build -o bin/app cmd/server/main.go

build-image:
	@docker build -t kf-usdtpay .

reload:
	@docker compose stop && docker compose rm -f && docker compose up -d
