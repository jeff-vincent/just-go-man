FROM golang:1.18

WORKDIR /usr/src/app

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]