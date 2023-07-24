FROM golang:1.18

WORKDIR /go/src/JarvisGo

COPY . .

# Use this only if your area is China mainland
RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go build -o server ./cmd/server/main.go

EXPOSE 5700 5701

ENTRYPOINT [ "./server" ]

# host ip record
# Linux
# 172.17.0.1
# Windows and MacOS
# host.docker.internal
