FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor

COPY . .

RUN go build -mod=vendor -o library cmd/main.go

EXPOSE 8080

CMD [ "./library" ]