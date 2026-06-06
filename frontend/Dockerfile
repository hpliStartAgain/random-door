# frontend/Dockerfile — 构建 React 产物
# ---------- 构建阶段 ----------
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./

# 配置国内 npm 镜像加速并安装依赖
RUN npm config set registry https://registry.npmmirror.com && npm ci --legacy-peer-deps

COPY . .
ARG VITE_AMAP_KEY
ARG VITE_AMAP_SECURITY_CODE
ARG VITE_API_BASE=/api
ENV VITE_AMAP_KEY=$VITE_AMAP_KEY
ENV VITE_AMAP_SECURITY_CODE=$VITE_AMAP_SECURITY_CODE
ENV VITE_API_BASE=$VITE_API_BASE
# VITE_ 前缀变量需在构建期注入。
RUN npm run build

# ---------- 运行阶段（用 caddy 托管 dist） ----------
FROM caddy:2-alpine
COPY --from=builder /app/dist /srv
# Caddyfile 由 compose 挂载，或在此 COPY 一个最简配置
EXPOSE 80
