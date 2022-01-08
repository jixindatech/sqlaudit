COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "GitNotFound")
BUILD_DATE=$(shell date '+%Y-%m-%d %H:%M:%S')
PWD=$(shell pwd)
all: build

build: sqlaudit
web:
	cd dashboard && yarn config set ignore-engines true && yarn install && npm run build:prod
sqlaudit: web
	cd $(PWD)
	go install golang.org/x/tools/cmd/goyacc@latest
	goyacc -o ./sqlparser/sql.go ./sqlparser/sql.y
	gofmt -w ./sqlparser/sql.go
	go build -ldflags "-X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\"" -o ./bin/sqlaudit ./cmd/sqlaudit
clean:
	@rm -rf bin
