## GoWebUser 后台管理系统

一个基于 Go 的轻量级后台管理系统，提供登录、注册、首页概览与用户管理等基础能力，采用分层架构，前端以模板页配合 Tailwind 与 Feather Icons 实现常见后台交互。

### 功能概览

- 登录 / 注册（表单校验、统一弹窗反馈）
- 首页概览（卡片统计、访问趋势占位）
- 用户管理（搜索、状态筛选、分页、状态/角色标签、删除）
- 角色可见性（管理员可见“新建/编辑/删除”入口，普通用户隐藏）

### 技术栈

- 后端：Go（net/http）、MySQL
  - 架构：Controller / Service / Repository / Infrastructure / Routes
  - 路由：http.ServeMux + 自定义注册
- 前端：模板页面（template/）
  - 样式：Tailwind CSS（CDN）
  - 图标：Feather Icons（CDN）
  - 组件：Bootstrap（主要用于 Modal 弹窗）
- 静态资源：static/css、static/js

### 目录结构（节选）

```
GoWebUser/
├─ main.go
├─ controller/         # 控制器
├─ services/           # 业务逻辑
├─ repositories/       # 数据访问
│  └─ repo_mysql/
├─ infrastructure/     # 基础设施（DB 等）
├─ routes/             # 路由注册
├─ template/           # 页面模板（login、register、index、userList）
└─ static/
   ├─ css/
   │  └─ common.css
   └─ js/
      ├─ common.js            # 通用弹窗
      ├─ login.js             # 登录提交逻辑
      ├─ register.js          # 注册提交逻辑
      ├─ userList.js          # 用户列表筛选/分页
      └─ login_register.js    # 登录/注册共享的“密码眼睛”逻辑
```

### 本地运行

1. 准备环境

- 安装 Go（建议 1.20+）与 MySQL
- 在 `infrastructure.InitDB()` 中配置数据库连接

1. 启动服务

```bash
go run .
```

访问页面：

- <http://localhost:8080/login>
- <http://localhost:8080/register>
- <http://localhost:8080/index>
- <http://localhost:8080/userList>

停止服务：在运行终端按 `Ctrl + C`。

### 接口约定（示例）

- POST `/api/login`         登录
- POST `/api/register`      注册
- GET  `/api/users/me`      获取当前登录用户信息（用于角色可见性）
- GET  `/api/users`         用户分页查询（支持 page/size/status/keyword）
- DELETE `/api/users/{id}`  删除用户（需管理员）

说明：以上为前端脚本使用到的典型接口，实际以后端实现为准；鉴权接口请携带 Cookie 等凭证，并在后端配置同域或正确的 CORS。