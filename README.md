# Galaxy Bing API

必应每日壁纸 API 服务，支持多地区、多尺寸的壁纸获取。

## 功能特性

- 支持多个地区的必应壁纸（zh-CN, de-DE, en-CA, en-GB, en-IN, en-US, fr-FR, it-IT, ja-JP）
- 提供今日壁纸、随机壁纸和历史壁纸列表
- 支持自定义图片尺寸（默认 1920x1080）
- 支持 JSON 和图片直接返回
- 自动同步最新壁纸（通过 GitHub Actions）
- 支持 API 访问控制

## 快速开始

1. 克隆项目

```bash
git clone https://github.com/gclm/galaxy-bing-api.git
cd galaxy-bing-api
```

2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，设置必要的配置
```

3. 初始化数据库

```bash
go run cmd/init/main.go
```

4. 启动服务

```bash
# 本地开发
go run cmd/server/main.go

# 或者使用 Vercel 开发环境
vercel dev
```

## API 文档

### 1. 获取今日壁纸

```http
GET /api/v1/today
```

查询参数：
- `mkt`: 地区代码，默认 zh-CN
- `w`: 图片宽度，默认 1920
- `h`: 图片高度，默认 1080
- `type`: 返回类型（image/json），默认 image

示例：
```bash
# 获取图片
curl "http://localhost:8080/api/v1/today?mkt=zh-CN"

# 获取 JSON
curl "http://localhost:8080/api/v1/today?type=json"
```

### 2. 获取随机壁纸

```http
GET /api/v1/random
```

参数与 today 接口相同。

### 3. 获取壁纸列表

```http
GET /api/v1/list
```

请求头：
- `Authorization`: API Token

查询参数：
- `page`: 页码，默认 1
- `pageSize`: 每页数量，默认 20
- `mkt`: 地区代码，可选

### 4. 获取指定日期壁纸

```http
GET /api/v1/date/{date}
```

请求头：
- `Authorization`: API Token

路径参数：
- `date`: 日期，格式：YYYY-MM-DD

查询参数：
- `mkt`: 地区代码，可选
- `w`: 图片宽度，默认 1920
- `h`: 图片高度，默认 1080
- `type`: 返回类型（image/json），默认 image

示例：
```bash
# 获取指定日期的图片
curl -H "Authorization: your-secret-token" "http://localhost:8080/api/v1/date/2024-02-19?mkt=zh-CN"

# 获取 JSON 格式
curl -H "Authorization: your-secret-token" "http://localhost:8080/api/v1/date/2024-02-19?type=json"
```

## 环境变量说明

```env
# MongoDB 配置
MONGODB_URI=mongodb+srv://<username>:<password>@<cluster>.mongodb.net/bing

# API 配置
PORT=8080
GIN_MODE=release
API_TOKEN=your-secret-token  # API 访问令牌

# 其他配置
PAGE_SIZE=20  # 默认分页大小
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t galaxy-bing-api .

# 运行容器
docker run -d \
  --name galaxy-bing-api \
  -p 8080:8080 \
  -e MONGODB_URI=your-mongodb-uri \
  -e API_TOKEN=your-secret-token \
  galaxy-bing-api
```

### 手动部署

1. 编译
```bash
go build -o bin/galaxy-bing-api
```

2. 运行
```bash
./bin/galaxy-bing-api
```

## 开发

### 本地开发

1. 安装依赖
```bash
npm install -g vercel
vercel login
```

2. 开发模式
```bash
# 使用 Go 开发服务器
go run cmd/server/main.go

# 或使用 Vercel 开发服务器
vercel dev --local-config=vercel.dev.json
```

3. 调试技巧
- 使用 `GIN_MODE=debug` 查看详细日志
- 修改 `vercel.dev.json` 配置本地环境变量
- 使用 `vercel logs` 查看部署日志

### 目录结构

```
├── cmd/                # 命令行工具
│   ├── fetch/         # 数据同步工具
│   └── init/          # 数据初始化工具
├── docs/              # 文档
├── internal/          # 内部包
│   ├── config/        # 配置管理
│   ├── database/      # 数据库操作
│   ├── handler/       # API 处理器
│   ├── logger/        # 日志管理
│   ├── middleware/    # 中间件
│   ├── model/         # 数据模型
│   └── utils/         # 工具函数
└── main.go            # 主程序
```

### 开发规范

- 遵循 Go 标准项目布局
- 使用 gofmt 格式化代码
- 添加必要的注释和文档
- 编写单元测试

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 Apache License 2.0 许可证。
