.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/podApi/podApi.proto

.PHONY: build
build:
	CGO_ENABLED="0" GOOS=linux GOARCH=amd64 go build -o podApi *.go

.PHONY: docker-build
docker-build:
	docker build -t baciyou/pod-api:latest .

.PHONY: docker-stop
docker-stop:
	docker stop gopass-pod-api

.PHONY: docker-rm
docker-rm:
	docker rm gopass-pod-api

.PHONY: docker-rmi
docker-rmi:
	docker image rm baciyou/pod-api:latest

.PHONY: docker-run
docker-run:
	docker run -d --name gopass-pod-api -p 8081:8081 -p 9092:9092 -p 9192:9192 -v /home/go-pro/src/podApi/micro.log:/micro.log baciyou/pod-api