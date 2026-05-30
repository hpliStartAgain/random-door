# Makefile — 常用命令封装，降低团队上手成本
# 对应 README.md。用法：make up / make seed ...

.PHONY: up down build logs migrate seed lint test fe-dev fe-build

# ---------- Docker ----------
up:            ## 启动全部容器
	docker compose up -d

down:          ## 停止并移除容器
	docker compose down

build:         ## 构建镜像
	docker compose build

logs:          ## 查看后端日志
	docker compose logs -f app

# ---------- 数据 ----------
# MySQL 已抽离为外部实例，自动建表与自动导入（Seed）由后端启动时自动完成。
migrate:       ## 已废弃，应用启动自动执行
	@echo "Migration is auto-executed on backend startup."

seed:          ## 已废弃，应用启动自动执行
	@echo "Seeding is auto-executed on backend startup."

# ---------- 质量 ----------
lint:          ## 后端 + 前端静态检查
	cd backend && gofmt -l . && go vet ./...
	cd frontend && npm run lint

test:          ## 后端单元测试（geo / achievement 为重点）
	cd backend && go test ./...

# ---------- 前端本地 ----------
fe-dev:        ## 前端开发服务器
	cd frontend && npm run dev

fe-build:      ## 前端构建产物
	cd frontend && npm run build