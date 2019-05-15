FROM golang:1.12-alpine

ENV \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN \
    apk add --no-cache --update tzdata git bash curl && \
    cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
    rm -rf /var/cache/apk/*

COPY . app
WORKDIR app

RUN \
    go version && \
    go mod download && \
    go build -ldflags "-X 'main.revision=$(git describe --abbrev=0 --tags 2> /dev/null || echo develop)'"

HEALTHCHECK --interval=15s --timeout=3s CMD curl --fail http://localhost:8080/health/ping || exit 1

EXPOSE 8080

ENTRYPOINT ["./pinky"]
