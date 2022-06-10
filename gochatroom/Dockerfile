FROM golang:1.17-stretch AS builder

WORKDIR /build
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gochatroomsvr

FROM debian:stretch-slim
WORKDIR /app

RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

COPY --from=builder /build/gochatroomsvr /app/gochatroomsvr
COPY ./wait-for.sh /app/wait-for.sh
COPY ./conf /app/conf
COPY .env.back /app/.env

RUN chmod +x /app/wait-for.sh

EXPOSE 10086
CMD ["/app/wait-for.sh", "mysql:3306", "--", "/app/gochatroomsvr", "-c", "conf/config.yaml"]
