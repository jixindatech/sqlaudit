COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "GitNotFound")
BUILD_DATE=$(shell date '+%Y-%m-%d %H:%M:%S')
PWD=$(shell pwd)
all: build

build: sqlaudit
web:
	cd dashboard && yarn install && npm run build:prod
sqlaudit:web
	cd $(PWD)
	go build -ldflags "-X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\"" -o ./bin/sqlaudit ./cmd/sqlaudit
clean:
	@rm -rf bin
