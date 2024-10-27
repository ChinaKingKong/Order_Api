# 订单系统 API

基于 Golang + Redis 实现的二级缓存订单系统，采用现代化的技术栈和最佳实践。

## 系统特点

- 二级缓存架构（本地缓存 + Redis）
- JWT 身份认证
- RESTful API 设计
- 统一错误处理和响应格式
- 完整的数据验证（支持中文错误提示）
- 全面的日志记录
- 优雅的项目结构

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **缓存**: Redis + 本地缓存
- **数据库**: MySQL
- **认证**: JWT
- **验证器**: validator/v10 + 中文翻译器

## 项目结构

```
order_api/
├── app/              # 应用程序核心
├── auth/             # JWT认证实现  
├── cache/            # 缓存层实现
├── config/           # 配置管理
├── database/         # 数据库管理
├── handler/          # HTTP处理器
├── middleware/       # 中间件
├── model/            # 数据模型
├── repository/       # 数据访问层
├── router/           # 路由配置
└── service/          # 业务逻辑层
```

## 主要功能

### 订单管理
- 创建订单
- 查询订单详情
- 更新订单状态
- 删除订单
- 订单列表查询

### 缓存设计
- 本地缓存（sync.Map）提供快速访问
- Redis 缓存提供分布式支持
- 缓存自动过期和更新机制
- 缓存一致性保证

### 认证授权
- JWT 令牌生成与验证
- 基于中间件的认证机制
- 令牌自动续期
- 用户角色控制

### 数据验证
- 请求参数自动校验
- 自定义验证规则
- 中文错误提示
- 统一的错误处理

## API 接口

### 认证接口
```
POST /api/v1/auth/login    # 用户登录
```

### 订单接口
```
POST   /api/v1/orders      # 创建订单
GET    /api/v1/orders/:id  # 获取订单详情
PUT    /api/v1/orders/:id  # 更新订单状态
DELETE /api/v1/orders/:id  # 删除订单
```

## 快速开始

1. 环境要求
   - Go 1.21+
   - MySQL 5.7+
   - Redis 6.0+

2. 配置文件
   - 复制 `config.json.example` 为 `config.json`
   - 修改数据库和 Redis 配置

3. 启动服务
```bash
# 安装依赖
go mod download

# 运行服务
go run main.go
```

## 配置说明

配置文件 `config.json` 包含以下主要配置：

```json
{
    "server": {
        "port": "8080",
        "read_timeout": 60,
        "write_timeout": 60
    },
    "database": {
        "host": "localhost",
        "port": "3306",
        "user": "root",
        "password": "root",
        "dbname": "orders"
    },
    "redis": {
        "host": "localhost",
        "port": "6379",
        "password": "",
        "db": 0
    },
    "jwt": {
        "secret_key": "your-secret-key",
        "token_expiry_hours": 24
    }
}
```

## 性能优化

1. 缓存策略
   - 二级缓存减少数据库访问
   - 缓存预热机制
   - 定时清理过期缓存

2. 数据库优化
   - 连接池管理
   - 索引优化
   - 软删除支持

3. 并发处理
   - goroutine 池
   - 并发安全的缓存访问
   - 分布式锁支持

## 部署建议

1. 安全配置
   - 修改默认的 JWT 密钥
   - 启用 HTTPS
   - 设置适当的访问控制

2. 监控告警
   - 系统资源监控
   - 接口性能监控
   - 错误日志监控

3. 高可用配置
   - 负载均衡
   - 服务器集群
   - 数据库主从复制

## 注意事项

1. 安全性
   - 生产环境必须修改默认的 JWT 密钥
   - 建议启用 HTTPS
   - 定期更新依赖包

2. 性能
   - 监控缓存命中率
   - 定期清理过期数据
   - 注意并发访问控制

3. 维护
   - 定期备份数据
   - 监控系统资源
   - 及时处理错误日志