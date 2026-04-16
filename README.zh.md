[English](./README.md) | [한국어](./README.ko.md) | [日本語](./README.ja.md) | **中文**

# Monkey Planner

> AI 智能体任务记忆仓库 — 类 Notion/JIRA 的问题追踪器 + MCP 服务器

一款协作工具，由人工创建并审批任务，AI 智能体通过 MCP（模型上下文协议）客户端来消费和执行任务。

![Monkey Planner](./docs/screenshots/d-home-l.png)

## 核心功能

### 问题与看板管理
- **看板视图** — 支持拖拽、横向滚动、筛选、排序和表格视图切换
- **创建问题** — 支持标题、Markdown 正文及自定义属性
- **自定义属性** — 支持 6 种类型
  - 文本 (text)
  - 数字 (number)
  - 单选 (select)
  - 多选 (multi_select)
  - 日期 (date)
  - 复选框 (checkbox)

### 审批流程
- **Pending → Approved** 专用审批端点（普通 PATCH 请求无法完成此操作）
- **审批队列** — 一键批量审批全部看板中的待审问题
- **Approved → InProgress → Done** — 状态可自由流转
- **Rejected 状态** — 可记录拒绝原因

### AI 智能体功能
- **智能体指令字段** — 为 MCP 智能体填写具体执行指令
- **验收标准** — 以清单形式管理完成条件
- **评论** — 记录每个问题的进展并进行沟通
- **依赖关系** — 表达问题之间的阻塞关系

### 数据可视化
- **日历** — 月度网格视图 + 每日统计（创建、审批、完成数量）
- **仪表盘** — 统计卡片 + 周活动图表
- **侧边栏** — 看板列表、问题计数、最近记录

### 用户体验
- **全局搜索** — 使用 Cmd+K 快速搜索
- **键盘快捷键**
  - `h` — 跳转到仪表盘
  - `a` — 跳转到审批队列
  - `?` — 查看快捷键帮助
  - `Cmd+S` — 保存
  - `Escape` — 关闭弹窗/对话框
- **侧边栏折叠/展开** — 优化屏幕空间利用
- **深色模式** — 主题切换
- **多语言** — 支持韩语、英语、日语、中文

### 自动化与集成
- **Webhook** — 支持 Discord、Slack、Telegram
  - 事件类型：`issue.created`、`issue.approved`、`issue.status_changed`、`issue.deleted`
- **JSON 导出** — 导出所有问题数据
- **右键上下文菜单** — 快捷操作菜单
- **问题模板** — 按看板存储至 localStorage

### MCP 服务器（AI 智能体集成）
10 个工具，赋能 AI 智能体自动化：
1. `list_boards` — 查看所有看板
2. `list_issues` — 查询问题（支持 boardId、status 筛选）
3. `get_issue` — 获取问题详情（含指令、标准、评论）
4. `create_issue` — 创建新问题
5. `approve_issue` — 审批：Pending → Approved
6. `claim_issue` — 认领：Approved → InProgress
7. `complete_issue` — 完成：InProgress → Done（可附带评论）
8. `add_comment` — 为问题添加评论
9. `update_criteria` — 勾选/取消验收标准
10. `search_issues` — 按标题搜索问题

## 技术栈

### 后端
- **语言**: Go 1.26
- **路由**: chi/v5
- **数据库**: SQLite / PostgreSQL（可选）
- **迁移**: goose/v3
- **文件内嵌**: embed.FS（单二进制文件部署）

### 前端
- **框架**: React 18
- **语言**: TypeScript
- **构建工具**: Vite 6
- **CSS**: Tailwind CSS
- **状态管理**: React Query (TanStack)
- **拖拽**: @dnd-kit/core, @dnd-kit/sortable
- **图标**: lucide-react
- **图表**: recharts
- **国际化**: react-i18next
- **Markdown**: react-markdown + rehype-sanitize

### MCP
- 协议: JSON-RPC 2.0 over stdio
- 适配: Claude Code、Claude Desktop

## 快速开始

### 环境要求
- Go 1.26 及以上
- Node.js 18 及以上
- npm 或 yarn

### 安装与运行

#### 1. 克隆仓库并初始化
```bash
git clone https://github.com/kjm99d/monkey-planner.git
cd monkey-planner
make init
```

#### 2. 生产构建（单二进制文件）
```bash
make build
./bin/monkey-planner
```

服务器将在 `http://localhost:8080` 启动，前端已内嵌其中。

#### 3. 开发模式（分离运行）

终端 1 — 后端：
```bash
make run-backend
```

终端 2 — 前端（Vite 开发服务器，端口 :5173）：
```bash
make run-frontend
```

前端会自动将 `/api` 请求代理到 `:8080`。

### 环境变量

```bash
# 服务器地址（默认: :8080）
export MP_ADDR=":8080"

# 数据库连接字符串
# SQLite（默认: sqlite://./data/monkey.db）
export MP_DSN="sqlite://./data/monkey.db"

# PostgreSQL 示例
export MP_DSN="postgres://user:password@localhost:5432/monkey_planner"
```

## MCP 服务器配置

### 快速设置（自动更新）

1. 从 [Releases](https://github.com/kjm99d/MonkeyPlanner/releases) 下载最新二进制文件（例如 `D:/mp/`）
2. 将 `update-and-run.sh` 下载到同一目录
3. 添加到 `.mcp.json` 或 Claude Code 设置中:

```json
{
  "mcpServers": {
    "monkey-planner": {
      "command": "bash",
      "args": ["/path/to/update-and-run.sh", "mcp"],
      "env": {
        "MP_DSN": "sqlite:///path/to/data/monkey.db"
      }
    }
  }
}
```

**Windows（无 bash）**: 使用 `update-and-run.bat`:

```json
{
  "mcpServers": {
    "monkey-planner": {
      "command": "D:/mp/update-and-run.bat",
      "args": ["mcp"],
      "env": {
        "MP_DSN": "sqlite://D:/mp/data/monkey.db"
      }
    }
  }
}
```

包装脚本会在每次启动时自动检查 GitHub 最新版本并更新二进制文件。

### 手动设置

不使用自动更新，直接配置:

```json
{
  "mcpServers": {
    "monkey-planner": {
      "command": "/path/to/monkey-planner",
      "args": ["mcp"],
      "env": {
        "MP_DSN": "sqlite:///path/to/data/monkey.db"
      }
    }
  }
}
```

Claude Code（`.mcp.json`）和 Claude Desktop（`~/.claude/claude_desktop_config.json`）使用相同格式。更改后请重启。

### MCP 工具使用示例

```
AI: 列出所有看板
→ 调用 list_boards()

AI: 查找与"认证"相关的问题
→ 调用 search_issues(query="认证")

AI: 审批第一个待处理问题，完成后提交 QA
→ 依次调用 approve_issue() → claim_issue() → submit_qa()
```

## 工作流 — 实际使用场景

以下是修复多语言切换 Bug 时的真实工作流，展示人类与 AI 智能体如何通过 Monkey Planner 协作。

### 状态流程

```
待处理 → 已审批 → 进行中 → QA验证 → 已完成
                    ↑              │（附带原因拒绝）
                    └──────────────┘
```

### 逐步说明

**1. 创建问题** — 人类发现 Bug，要求 AI 注册问题
```
人类: "点击语言切换按钮没有弹出下拉菜单，创建一个问题。"
AI:   create_issue(boardId, title, body, instructions)  →  状态: 待处理
```

**2. 审批** — 人类确认并审批
```
人类: （在看板上点击 Approve，或指示 AI）
AI:   approve_issue(issueId)  →  状态: 已审批
```

**3. 开始工作** — AI 认领问题并开始修复
```
AI:   claim_issue(issueId)  →  状态: 进行中
      - 分析代码，定位根因
      - 实现修复，运行测试
      - 提交更改
```

**4. 提交 QA** — 完成工作后请求验证
```
AI:   submit_qa(issueId, comment: "提交 abc1234 — 修复点击处理器")
      →  状态: QA验证
      add_comment(issueId, "提交信息: ...")
```

**5. 验证** — 人类亲自测试
```
人类: 在浏览器中测试，发现下拉菜单被侧边栏遮挡
      →  reject_issue(issueId, reason: "下拉菜单被侧边栏遮挡")
      →  状态: 进行中（回到第3步）

      或者

人类: 重新修复后测试，一切正常
      →  complete_issue(issueId)  →  状态: 已完成
```

**6. 反馈循环** — 通过评论沟通
```
人类: add_comment("下拉菜单被左侧遮挡了，请修复")
AI:   get_issue() → 查看评论 → 修复 → 提交 → submit_qa()
人类: 测试 → complete_issue()  →  完成 ✓
```

### 关键要点

- **人类控制关卡**: 审批、QA 通过/拒绝、完成决定
- **AI 执行工作**: 代码分析、实现、测试、提交
- **评论是沟通渠道**: 双方通过 `add_comment` 交换反馈
- **QA 循环防止过早完成**: 必须通过人类验证才能标记为完成

## API 文档

OpenAPI 3.0 规范：[backend/docs/swagger.yaml](./backend/docs/swagger.yaml)

### 主要端点

#### Boards（看板）
```
GET    /api/boards                  # 获取看板列表
POST   /api/boards                  # 创建看板
PATCH  /api/boards/{id}             # 修改看板
DELETE /api/boards/{id}             # 删除看板
```

#### Issues（问题）
```
GET    /api/issues                  # 获取问题列表（筛选: boardId, status, parentId）
POST   /api/issues                  # 创建问题
GET    /api/issues/{id}             # 获取问题详情（含子问题）
PATCH  /api/issues/{id}             # 修改问题（状态、属性、标题等）
DELETE /api/issues/{id}             # 删除问题
POST   /api/issues/{id}/approve     # 审批问题（Pending → Approved）
```

#### Comments（评论）
```
GET    /api/issues/{issueId}/comments    # 获取评论列表
POST   /api/issues/{issueId}/comments    # 添加评论
DELETE /api/comments/{commentId}         # 删除评论
```

#### Properties（自定义属性）
```
GET    /api/boards/{boardId}/properties      # 获取属性定义列表
POST   /api/boards/{boardId}/properties      # 创建属性
PATCH  /api/boards/{boardId}/properties/{propId}  # 修改属性
DELETE /api/boards/{boardId}/properties/{propId}  # 删除属性
```

#### Webhooks
```
GET    /api/boards/{boardId}/webhooks       # 获取 Webhook 列表
POST   /api/boards/{boardId}/webhooks       # 创建 Webhook
PATCH  /api/boards/{boardId}/webhooks/{whId}    # 修改 Webhook
DELETE /api/boards/{boardId}/webhooks/{whId}    # 删除 Webhook
```

#### Calendar（日历）
```
GET /api/calendar           # 月度统计（year、month 必填）
GET /api/calendar/day       # 当日问题列表（date 必填）
```

完整 Schema 请参阅 [backend/docs/swagger.yaml](./backend/docs/swagger.yaml)。

## 项目结构

```
monkey-planner/
├── backend/
│   ├── cmd/monkey-planner/
│   │   ├── main.go              # 入口（HTTP 服务器）
│   │   └── mcp.go               # MCP 服务器（JSON-RPC stdio）
│   ├── internal/
│   │   ├── domain/              # 领域模型（Issue、Board 等）
│   │   ├── service/             # 业务逻辑
│   │   ├── storage/             # 数据库（SQLite/PostgreSQL）
│   │   ├── http/                # HTTP 处理器 & 路由
│   │   └── migrations/          # goose 迁移文件
│   ├── web/                     # 前端内嵌（embed.FS）
│   ├── docs/
│   │   └── swagger.yaml         # OpenAPI 3.0 规范
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── components/          # 通用组件
│   │   ├── features/            # 页面 & 功能组件
│   │   │   ├── home/           # 仪表盘
│   │   │   ├── board/          # 看板
│   │   │   ├── issue/          # 问题详情
│   │   │   ├── calendar/       # 日历
│   │   │   └── approval/       # 审批队列
│   │   ├── api/                 # API Hooks & 客户端
│   │   ├── design/              # Tailwind 设计令牌
│   │   ├── i18n/                # 多语言（en.json, ko.json, ja.json, zh.json）
│   │   ├── App.tsx              # 路由
│   │   ├── index.css            # 全局样式
│   │   └── main.tsx
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── tailwind.config.js
│
├── .mcp.json                    # Claude Code MCP 配置
├── Makefile                     # 构建 & 开发命令
├── .githooks/                   # Git Hooks
└── data/                        # SQLite 数据库（默认）
```

## 测试

### 后端测试
```bash
make test-backend
```

### 前端测试
```bash
make test-frontend
```

### 无障碍测试
```bash
make test-a11y
```

### 全量测试
```bash
make test
```

## 常用命令

```bash
# 克隆后初始化
make init

# 生产构建
make build

# 启动生产服务器
./bin/monkey-planner

# 开发模式
make run-backend        # 终端 1
make run-frontend       # 终端 2

# 清理构建产物
make clean
```

## 状态流转规则

```
Pending
  ↓ (approve endpoint)
Approved
  ↓ (PATCH status)
InProgress
  ↓ (PATCH status)
Done

Pending → Approved: POST /api/issues/{id}/approve（专用端点）
Approved ↔ InProgress ↔ Done: 通过 PATCH 自由流转
Pending: 一旦离开，无法从其他状态回退
Rejected: 拒绝状态（独立追踪）
```

## 许可证

MIT

## 贡献

欢迎提交 Issue 和 Pull Request。

## 联系我们

如有任何问题或反馈，请通过 GitHub Issues 提交。
