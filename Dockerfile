FROM golang:alpine

WORKDIR /go/src
COPY .dockerignore /go/src/.dockerignore
ADD . /go/src
RUN go env GO111MODULE=on
RUN cd /go/src && go mod tidy && go mod download
RUN go build -o server

RUN chmod +x ./server
EXPOSE 50051
CMD ["./server"]
