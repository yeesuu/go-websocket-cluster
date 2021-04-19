FROM golang:1.16 AS builder
ENV GO111MODULE=on
ENV GOPROXY='https://goproxy.cn'
WORKDIR /var/www
COPY .. .
RUN go build

FROM centos:latest as server
WORKDIR /var/www
COPY --from=builder /var/www/go-websocket-cluster ./
COPY --from=builder /var/www/config.example.yml ./config.yml
CMD ["./go-websocket-cluster"]
EXPOSE 8000