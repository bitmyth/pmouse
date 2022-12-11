FROM golang:1.17-alpine as build-stage

RUN apk --update upgrade \
    && apk --no-cache --no-progress add git mercurial bash gcc musl-dev curl tar ca-certificates tzdata \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.mod .

RUN GO111MODULE=on GOPROXY=https://goproxy.cn go mod download

COPY . .

RUN GOGC=off  go build -v -o dist/app .


FROM scratch
COPY --from=build-stage /app/dist/app /
WORKDIR /
VOLUME /download
EXPOSE 8000
ENTRYPOINT ["/app"]