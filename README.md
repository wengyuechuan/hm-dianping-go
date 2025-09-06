# 黑马点评 Go 版本

**⚠️ 学习项目声明：本项目是黑马点评的 Go 语言学习版本，仅供教育和学习目的使用，不得用于商业用途。**

基于 Go 语言和 Gin 框架的黑马点评后端项目，采用分层架构设计。本项目参考原 Java 版本的黑马点评项目，使用 Go 语言重新实现，旨在帮助开发者学习 Go 语言的 Web 开发技术栈。

## 项目结构

```
hm-dianping-go/
├── config/                 # 配置层
│   ├── config.go          # 配置管理
│   └── application.yaml   # 配置文件
├── dao/                   # 数据访问层
│   └── db.go             # 数据库连接管理
├── handler/               # 处理器层
│   ├── user.go           # 用户相关处理器
│   └── common.go         # 通用处理器
├── router/                # 路由层
│   └── router.go         # 路由配置
├── service/               # 服务层
│   ├── models.go         # 数据模型
│   ├── user.go           # 用户服务
│   └── common.go         # 通用服务
├── utils/                 # 工具层
│   ├── response.go       # 响应工具
│   ├── jwt.go            # JWT工具
│   ├── password.go       # 密码工具
│   └── middleware.go     # 中间件
├── go.mod                # Go模块文件
├── main.go               # 项目入口
└── README.md             # 项目说明
```

## 技术栈

- **Web框架**: Gin
- **数据库**: MySQL + GORM
- **缓存**: Redis
- **认证**: JWT
- **配置**: YAML
- **密码加密**: bcrypt

## 功能模块

### 用户模块
- 用户注册
- 用户登录
- 获取用户信息
- 更新用户信息

### 商铺模块
- 获取商铺列表
- 根据ID获取商铺
- 根据类型获取商铺

### 优惠券模块
- 获取优惠券列表
- 秒杀优惠券

### 博客模块
- 获取博客列表
- 根据ID获取博客
- 创建博客
- 点赞博客

### 关注模块
- 关注用户
- 取消关注
- 获取共同关注

## 快速开始

### 1. 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置数据库

修改 `config/application.yaml` 文件中的数据库配置：

```yaml
database:
  host: "localhost"
  port: "3306"
  username: "root"
  password: "your_password"
  dbname: "hmdp"
  charset: "utf8mb4"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0
```

### 4. 运行项目

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

## API 接口

### 用户接口

- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录
- `GET /api/user/info` - 获取用户信息（需要认证）
- `PUT /api/user/update` - 更新用户信息（需要认证）

### 商铺接口

- `GET /api/shop/list` - 获取商铺列表
- `GET /api/shop/:id` - 根据ID获取商铺
- `GET /api/shop/type/:typeId` - 根据类型获取商铺

### 优惠券接口

- `GET /api/voucher/list` - 获取优惠券列表
- `POST /api/voucher/seckill/:id` - 秒杀优惠券（需要认证）

### 博客接口

- `GET /api/blog/list` - 获取博客列表
- `GET /api/blog/:id` - 根据ID获取博客
- `POST /api/blog/create` - 创建博客（需要认证）
- `POST /api/blog/like/:id` - 点赞博客（需要认证）

### 关注接口

- `POST /api/follow/:id` - 关注用户（需要认证）
- `DELETE /api/follow/:id` - 取消关注（需要认证）
- `GET /api/follow/common/:id` - 获取共同关注（需要认证）

### 健康检查

- `GET /health` - 健康检查

## 项目特点

1. **分层架构**: 采用经典的分层架构，职责清晰，易于维护
2. **统一响应**: 统一的API响应格式，便于前端处理
3. **JWT认证**: 基于JWT的用户认证机制
4. **中间件支持**: 支持CORS、日志、认证等中间件
5. **配置管理**: 基于YAML的配置管理，支持不同环境配置
6. **数据库连接池**: 合理的数据库连接池配置
7. **错误处理**: 完善的错误处理机制

## 开发说明

**重要提醒：本项目仅为学习目的而创建，基于黑马程序员的教学项目进行 Go 语言实现。**

本项目参考Java版黑马点评项目架构，使用Go语言重新实现。项目采用模块化设计，各层职责明确：

- **Router层**: 负责路由配置和请求分发
- **Handler层**: 负责请求参数验证和响应处理
- **Service层**: 负责业务逻辑实现
- **DAO层**: 负责数据库操作
- **Config层**: 负责配置管理
- **Utils层**: 负责公共工具函数

## 注意事项

1. **本项目仅供学习使用，请勿用于生产环境或商业用途**
2. 请确保MySQL和Redis服务正常运行
3. 首次运行前请创建对应的数据库
4. 生产环境请修改JWT密钥和数据库密码
5. 建议使用GORM的自动迁移功能创建数据表
6. 支持通过命令行参数指定配置文件：`go run main.go -config=path/to/config.yaml`

## 许可证

MIT License