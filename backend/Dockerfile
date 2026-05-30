FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .
RUN go build -o server ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/data ./data
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations

RUN mkdir -p uploads

EXPOSE 8080
CMD ["./server"]
