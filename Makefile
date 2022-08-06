build:
	go build -o ./.bin main.go

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
	scp pi@192.168.50.55:/docker-vk/sqlite/stage.db .

#docker run -e TOKEN='vk1.a._wYEMXZSELBhBsjIyQfCS_J-l-mfmC68PbpDuVK-_-ObVLWrMGWr3JlJpCAr4r8liAOA7N4iRKHt_kg3qCe3QB3wjOq9u2QGAFdhsh3CcgplGbdi2zDZbHXdDA9s9CmwUW6D5L7yUA9EBSPV_4Gh_hccYNiJ7FvRFl4FAkNtDBJ6HUZS9SIrg1O6doqHSRxp' -e MY_ID='192398160' -d -v UPLOADer-db:/docker-vk/sqlite --restart always --publish 80:80 --name UPLOADer mrdjeb/vk:latest