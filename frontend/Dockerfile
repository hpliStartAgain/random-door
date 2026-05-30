# frontend/Dockerfile — 构建 React 产物（可选；MVP 也可本地 build 后交 Caddy）
# ---------- 构建阶段 ----------
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
# VITE_ 前缀变量需在构建期注入（通过 --build-arg 或 .env）
RUN npm run build

# ---------- 运行阶段（用 caddy 托管 dist） ----------
FROM caddy:2-alpine
COPY --from=builder /app/dist /srv
# Caddyfile 由 compose 挂载，或在此 COPY 一个最简配置
EXPOSE 80