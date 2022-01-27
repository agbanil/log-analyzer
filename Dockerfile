FROM golang:1.17.6

WORKDIR /app

COPY . .

RUN go get -d -v ./...

ENTRYPOINT ["go", "run", "main.go"]