FROM golang:1.22.2-alpine3.19 AS builder


WORKDIR /app

RUN apk update && apk add --no-cache git
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o datingapp ./cmd/app/main.go


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/datingapp .

EXPOSE 8080
CMD ["./datingapp"]
