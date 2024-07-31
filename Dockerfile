FROM golang:1.22.4-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /proxy-server ./cmd/app/main.go

CMD ["/proxy-server"]
