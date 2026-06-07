# 目录结构总览

## 1. 仓库根

```text
random-door/
├── AGENTS.md
├── CLAUDE.md
├── README.md
├── CHANGELOG.md
├── TODO.md
├── Makefile
├── .env.example
├── .gitignore
├── Caddyfile
├── docker-compose.yml
├── backend/
├── frontend/
├── docs/
└── scripts/
```

## 2. backend/

```text
backend/
├── cmd/
│   ├── server/main.go
│   └── seedtool/main.go
├── internal/
│   ├── api/          router + handlers
│   ├── service/      business orchestration
│   ├── repository/   GORM data access
│   ├── model/        GORM models
│   ├── geo/          random roaming algorithm
│   ├── ai/           LLM and image clients, prompt builder
│   ├── upload/       upload validation and storage
│   ├── seed/         seed validation/import
│   ├── achievement/  rules and evaluator
│   ├── middleware/   cors/logger/recover/rate limit
│   └── config/       env config
├── data/seed/        cities.json, achievements.json
├── migrations/       schema.sql
├── static/           checked-in static assets
├── Dockerfile
├── go.mod
└── go.sum
```

## 3. frontend/

```text
frontend/
├── public/           favicon/icon/fox images
├── src/
│   ├── App.tsx
│   ├── main.tsx
│   ├── api/          client, API functions, types
│   ├── assets/       asset mapping
│   ├── components/   map, sidebar, drawer, checkin, overlays
│   ├── lib/          city filters, share image
│   ├── pages/        admin, assets, achievements, guess challenge
│   └── store/        Zustand stores
├── Dockerfile
├── package.json
├── package-lock.json
├── vite.config.ts
└── tailwind.config.js
```

## 4. docs/

```text
docs/
├── agent/   coding agent rules
├── arch/    architecture and deployment docs
├── design/  API, database, backend, frontend, algorithm, AI, achievement details
└── product/ PRD, structure, flows, acceptance, demo script
```

## 5. scripts/

```text
scripts/
├── ai_smoke.py
├── build_seed.py
├── fetch_real_assets.py
├── requirements.txt
├── seed_inputs.csv
└── seed_builder/
```

## 6. 运行时目录

```text
uploads/
  selfies/
  scenes/
  generated/
  guess/
  admin_imports/
```

`uploads/`、构建产物、依赖目录、缓存文件和密钥文件不得提交。
