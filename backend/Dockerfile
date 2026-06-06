FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .
RUN mkdir -p static uploads
RUN go build -o server ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/data ./data
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/static ./static

# 自动创建 uploads 目录；static 已由镜像复制，避免 named volume 遮住内置演示素材
RUN mkdir -p uploads

EXPOSE 8080
CMD ["./server"]
