FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 as prod
WORKDIR /app
COPY --from=builder /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]