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
COPY --from=builder /app/migrations ./migrations

# 自动创建 uploads 和 static 目录，避免因目录为空未被 Git/Docker 复制而导致构建失败
RUN mkdir -p uploads static

EXPOSE 8080
CMD ["./server"]
