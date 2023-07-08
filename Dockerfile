# compile
FROM golang:1.20 AS buildStage
WORKDIR /go/src
COPY .dockerignore /go/src/.dockerignore
COPY . /go/src
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
RUN cd /go/src && go mod tidy && go mod download
RUN go build -o myapp

# pack
FROM alpine:latest
WORKDIR /app
COPY --from=buildStage /go/src/myapp /app/

# portf
EXPOSE 50051

# run
ENTRYPOINT /app/myapp