# Initial stage: download modules
FROM golang:1.18-alpine as builder

ENV config=docker

WORKDIR /app

COPY go.mod /app
COPY go.sum /app

COPY ./ /app

RUN go mod download


# Intermediate stage: Build the binary
FROM golang:1.18-alpine as runner

RUN apk add git

COPY --from=builder ./app ./app

RUN go env -w GO111MODULE=auto

RUN go get -u github.com/githubnemo/CompileDaemon

WORKDIR /app
ENV config=docker

EXPOSE 4000 

ENTRYPOINT CompileDaemon --build="go build cmd/api/main.go" --command=./main