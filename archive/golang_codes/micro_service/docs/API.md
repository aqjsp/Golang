# 分布式微服务架构平台 API 文档

## 概述

本文档描述了分布式微服务架构平台的API接口规范，包括API网关和各个微服务的接口定义。

## 基础信息

- **API版本**: v1.0.0
- **基础URL**: http://localhost:8080
- **认证方式**: JWT Bearer Token
- **数据格式**: JSON
- **字符编码**: UTF-8

## 通用响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "error message",
  "data": null
}
```

### 状态码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权，需要登录 |
| 403 | 禁止访问，权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## API网关接口

### 健康检查

**接口地址**: `GET /health`

**接口描述**: 检查API网关健康状态

**请求参数**: 无

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T12:00:00Z",
    "version": "v1.0.0"
  }
}
```

### 服务发现

**接口地址**: `GET /admin/services`

**接口描述**: 获取所有注册的服务列表

**请求参数**: 无

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "services": [
      {
        "name": "user-service",
        "address": "localhost",
        "port": 8081,
        "tags": ["user", "auth"],
        "health": true
      }
    ]
  }
}
```

### 指标监控

**接口地址**: `GET /metrics`

**接口描述**: 获取Prometheus监控指标

**请求参数**: 无

**响应格式**: Prometheus格式

## 用户服务接口

### 用户注册

**接口地址**: `POST /api/v1/users/register`

**接口描述**: 用户注册

**请求参数**:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "phone": "13800138000"
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名，3-50字符 |
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码，至少6位 |
| phone | string | 否 | 手机号码 |

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "phone": "13800138000",
    "status": 1,
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

### 用户登录

**接口地址**: `POST /api/v1/users/login`

**接口描述**: 用户登录

**请求参数**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名或邮箱 |
| password | string | 是 | 密码 |

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-02T12:00:00Z",
    "user": {
      "user_id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "phone": "13800138000",
      "status": 1
    }
  }
}
```

### 获取用户信息

**接口地址**: `GET /api/v1/users/{id}`

**接口描述**: 根据用户ID获取用户信息

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 用户ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "phone": "13800138000",
    "avatar": "",
    "status": 1,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

### 更新用户信息

**接口地址**: `PUT /api/v1/users/{id}`

**接口描述**: 更新用户信息

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 用户ID |

**请求参数**:
```json
{
  "email": "newemail@example.com",
  "phone": "13900139000",
  "avatar": "http://example.com/avatar.jpg"
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| email | string | 否 | 邮箱地址 |
| phone | string | 否 | 手机号码 |
| avatar | string | 否 | 头像URL |

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_id": 1,
    "username": "testuser",
    "email": "newemail@example.com",
    "phone": "13900139000",
    "avatar": "http://example.com/avatar.jpg",
    "status": 1,
    "updated_at": "2024-01-01T13:00:00Z"
  }
}
```

### 删除用户

**接口地址**: `DELETE /api/v1/users/{id}`

**接口描述**: 删除用户（软删除）

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 用户ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": null
}
```

### 用户列表

**接口地址**: `GET /api/v1/users`

**接口描述**: 获取用户列表（分页）

**请求头**:
```
Authorization: Bearer {token}
```

**查询参数**:
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 10 | 每页数量 |
| keyword | string | 否 | - | 搜索关键词 |
| status | int | 否 | - | 用户状态 |

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "users": [
      {
        "user_id": 1,
        "username": "testuser",
        "email": "test@example.com",
        "phone": "13800138000",
        "status": 1,
        "created_at": "2024-01-01T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 1,
      "pages": 1
    }
  }
}
```

## 订单服务接口

### 创建订单

**接口地址**: `POST /api/v1/orders`

**接口描述**: 创建新订单

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "items": [
    {
      "product_id": 1001,
      "product_name": "商品A",
      "product_price": 99.99,
      "quantity": 2
    }
  ],
  "shipping_address": "北京市朝阳区xxx街道xxx号",
  "remark": "请尽快发货"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "order_id": 1,
    "order_no": "ORD202401010001",
    "user_id": 1,
    "total_amount": 199.98,
    "status": 1,
    "payment_status": 0,
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

## 支付服务接口

### 创建支付

**接口地址**: `POST /api/v1/payments`

**接口描述**: 创建支付订单

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "order_no": "ORD202401010001",
  "amount": 199.98,
  "payment_method": 1
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| order_no | string | 是 | 订单号 |
| amount | decimal | 是 | 支付金额 |
| payment_method | int | 是 | 支付方式：1-支付宝，2-微信，3-银行卡 |

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "payment_id": 1,
    "payment_no": "PAY202401010001",
    "order_no": "ORD202401010001",
    "amount": 199.98,
    "payment_method": 1,
    "payment_status": 0,
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 10001 | 参数验证失败 |
| 10002 | 用户名已存在 |
| 10003 | 邮箱已存在 |
| 10004 | 用户不存在 |
| 10005 | 密码错误 |
| 10006 | 用户已被禁用 |
| 10007 | Token无效 |
| 10008 | Token已过期 |
| 20001 | 订单不存在 |
| 20002 | 订单状态错误 |
| 30001 | 支付失败 |
| 30002 | 支付金额错误 |
| 90001 | 服务不可用 |
| 90002 | 数据库连接失败 |
| 90003 | 缓存连接失败 |

## 认证说明

### JWT Token格式

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Token载荷

```json
{
  "user_id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "role": "user",
  "exp": 1704110400,
  "iat": 1704024000,
  "iss": "micro-service-platform",
  "sub": "testuser"
}
```

## 限流说明

- 全局限流：每秒1000次请求
- 单IP限流：每分钟100次请求
- 登录接口：每分钟5次请求
- 注册接口：每分钟3次请求

## 监控指标

系统提供以下监控指标：

- 请求总数
- 请求响应时间
- 错误率
- 活跃连接数
- 数据库连接池状态
- 缓存命中率
- 服务健康状态

监控数据可通过 `/metrics` 端点获取，格式为Prometheus标准格式。

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| v1.0.0 | 2024-01-01 | 初始版本，包含用户、订单、支付基础功能 |