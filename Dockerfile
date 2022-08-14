###     Need support arm64v8 architecture, case Rasbperry Pi 3 Model B works on this type.
###     Need support Debian GNU/Linux 11 (bullseye) OS-release, case Rasbperry Pi 3 Model B works on this type.

FROM golang:alpine AS builder

#RUN go version

COPY . /github.com/MrDjeb/vk/
WORKDIR /github.com/MrDjeb/vk/

RUN apk add build-base
#RUN go mod download 
RUN GOOS=linux GOARCH=arm64 go build -o ./.bin main.go

#FROM alpine:latest
FROM scratch as final

WORKDIR /docker-vk/

COPY --from=builder /github.com/MrDjeb/vk/.bin .
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
#COPY --from=0 /github.com/MrDjeb/vk/configs configs/

ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip
EXPOSE 80

ENTRYPOINT ["./.bin"]