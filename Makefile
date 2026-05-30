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
migrate:       ## 执行建表（若未在 mysql 初始化时自动执行）
	docker compose exec -T mysql sh -c 'mysql -u$(DB_USER) -p$(DB_PASSWORD) $(DB_NAME)' < backend/migrations/schema.sql

seed:          ## 导入 12 城 + 成就种子（由后端提供 seed 子命令或启动自动导入）
	docker compose exec app /app/server -seed

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