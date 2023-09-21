VERSION=v2.3.1
DATETIME=$(shell date "+%Y%m%d%H%M%S")

compile:
	@go mod tidy && go build -o ./bin/cm-api-server

docker:
	@docker build -t richzhao/cm-api-server:${VERSION} .
