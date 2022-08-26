build:
	go build -x -o ./.bin main.go

run: build
	./.bin

build-c:
	docker build  --platform linux/arm64/v8 -t docker-vk:v0 .

start-c:
	docker run --name UPLOADer --env-file .env -p 80:80 docker-vk:v0

push-c:
	docker tag docker-vk:v0 mrdjeb/vk
	docker push mrdjeb/vk

get-db:
	scp pi@raspi:/docker-vk/sqlite/stage.db .