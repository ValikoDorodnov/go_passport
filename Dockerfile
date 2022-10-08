FROM golang:1.18.2-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update

# build go app
RUN go mod download
RUN go build -o passport ./cmd/app/main.go

CMD ["./passport"]