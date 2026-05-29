# user-system-final

一个基于 Go + Gin + JWT + Redis 的轻量级用户系统示例项目，主要实现用户登录、退出登录、登录态校验、用户信息查询、Redis 缓存和接口限流。

## 项目简介

本项目是一个后端 API 服务 Demo，围绕用户认证和用户资料访问展开。项目使用 Gin 提供 HTTP 接口，使用 JWT 生成访问令牌，使用 Redis 保存登录态、缓存用户资料，并通过 Redis 计数实现基于用户维度的接口限流。

当前项目使用内存仓储保存测试用户，适合用于学习 Go Web 项目分层、JWT 认证、Redis 缓存、限流和缓存击穿处理等后端常见场景。

## 技术栈

- Go 1.25.4
- Gin：HTTP Web 框架
- JWT：用户 token 生成与解析
- Redis：登录态存储、用户信息缓存、接口限流计数
- go-redis：Redis 客户端
- singleflight：合并并发请求，降低缓存击穿风险

## 核心功能

- 用户登录：校验用户名和密码，生成 JWT token
- 登录态保存：将 token 写入 Redis，支持服务端主动失效
- 用户退出：删除 Redis 中的 token
- 认证中间件：校验 Authorization 请求头中的 Bearer token
- 用户资料查询：读取用户信息并使用 Redis 缓存
- 接口限流：对已登录用户按分钟限制访问次数
- 缓存击穿保护：使用 singleflight 合并同一资源的并发查询

## 项目结构

```text
.
├── cmd
│   └── main.go                  # 应用入口，初始化路由、服务和中间件
├── internal
│   ├── auth
│   │   ├── jwt.go               # JWT 生成与解析
│   │   ├── middleware.go        # 登录认证中间件
│   │   └── rate_limit.go        # Redis 限流中间件
│   ├── cache
│   │   └── redis.go             # Redis 客户端封装
│   ├── handler
│   │   └── user_handler.go      # HTTP 接口处理层
│   ├── logger
│   │   └── logger.go            # 简单日志封装
│   ├── model
│   │   └── user.go              # 用户模型
│   ├── repository
│   │   ├── memory_repo.go       # 内存用户仓储实现
│   │   └── user_repository.go   # 用户仓储接口
│   └── service
│       └── user_service.go      # 用户业务逻辑
├── go.mod
├── go.sum
└── test.http                    # 接口测试示例
```

## 快速开始

### 1. 启动 Redis

项目默认连接本地 Redis：

```text
localhost:6379
```

请先确保 Redis 已经启动。

### 2. 启动服务

```bash
go run ./cmd
```

服务默认监听：

```text
http://localhost:8080
```

## 测试账号

项目内置了一个测试用户：

```text
username: test
password: 123456
```

该用户定义在 `internal/repository/memory_repo.go` 中。

## API 接口

### 登录

```http
POST /login
Content-Type: application/json
```

请求体：

```json
{
  "username": "test",
  "password": "123456"
}
```

成功响应：

```json
{
  "token": "<jwt_token>"
}
```

### 获取用户信息

```http
GET /profile
Authorization: Bearer <jwt_token>
```

说明：

- 需要携带登录接口返回的 token
- 请求会先通过认证中间件
- 用户资料会优先从 Redis 缓存读取
- 缓存未命中时会进入业务层查询并重新写入缓存

### 退出登录

```http
POST /logout
Authorization: Bearer <jwt_token>
```

成功响应：

```json
{
  "msg": "logout success"
}
```

退出登录后，对应 token 会从 Redis 删除，再次访问受保护接口会认证失败。

## 认证与缓存流程

1. 客户端调用 `/login` 提交用户名和密码。
2. 服务端校验用户信息，生成 JWT token。
3. 服务端将 token 写入 Redis，key 格式为 `login:token:<token>`。
4. 客户端访问 `/profile` 或 `/logout` 时，在请求头中携带 `Authorization: Bearer <token>`。
5. 认证中间件先解析 JWT，再检查 Redis 中是否存在该 token。
6. 校验通过后，将 userID 写入 Gin context，交给后续业务处理。

## 限流设计

受保护接口挂载了 Redis 限流中间件，当前限制为每个用户每分钟最多 5 次请求。

限流 key 格式：

```text
rate_limit:<userID>:<yyyyMMddHHmm>
```

当请求次数超过限制时，接口返回：

```json
{
  "error": "rate limit exceeded"
}
```

## 项目亮点

- 分层结构清晰，便于后续扩展数据库、配置管理和更多业务接口
- JWT 与 Redis 结合，既保留无状态 token 的使用方式，又支持服务端主动登出
- 使用 Redis 同时完成 token 存储、用户缓存和访问限流
- 使用 singleflight 合并并发查询，降低缓存击穿带来的重复访问压力
- 缓存过期时间加入随机偏移，降低热点缓存同时失效的概率
- repository 使用接口抽象，内存实现可以平滑替换为数据库实现

## 后续可优化方向

- 将 Redis 地址、服务端口、JWT secret 改为配置文件或环境变量
- 使用数据库替换内存仓储
- 对密码进行哈希存储，例如 bcrypt
- 增加请求参数校验和统一响应结构
- 增加单元测试和接口测试
- 完善日志字段，支持结构化日志
- 增加注册、刷新 token、修改密码等用户功能
