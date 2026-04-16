[English](./README.md) | [한국어](./README.ko.md) | **日本語** | [中文](./README.zh.md)

# Monkey Planner

> AIエージェントのタスク管理ストア — Notion/JIRAスタイルのイシュートラッカー + MCPサーバー

人間がイシューを作成・承認し、AIエージェントがMCP（Model Context Protocol）クライアント経由でタスクを処理する協働ツールです。

![Monkey Planner](./docs/screenshots/d-home-l.png)

## 主な機能

### イシュー & ボード管理
- **カンバンボード** — ドラッグ&ドロップ、横スクロール、フィルタリング、ソート、テーブルビュー切り替え
- **イシュー作成** — タイトル、Markdownの本文、カスタム属性に対応
- **カスタム属性** — 6種類のタイプをサポート
  - テキスト (text)
  - 数値 (number)
  - 選択 (select)
  - 複数選択 (multi_select)
  - 日付 (date)
  - チェックボックス (checkbox)

### 承認フロー (Approval Flow)
- **Pending → Approved** 専用の承認エンドポイント（通常のPATCHでは変更不可）
- **承認キュー** — 全ボードのPendingイシューを一括承認
- **Approved → InProgress → Done** — 自由なステータス遷移
- **Rejected ステータス** — 却下理由の記録が可能

### エージェント機能
- **エージェント指示フィールド** — MCPエージェントが参照する具体的な指示を入力
- **成功基準** — 完了条件をチェックリストで管理
- **コメント** — イシューごとの進捗記録とコミュニケーション
- **依存関係** — イシュー間のブロック関係を表現

### データ可視化
- **カレンダー** — 月間グリッド + 日別実績（作成・承認・完了のカウント）
- **ダッシュボード** — 統計カード + 週次アクティビティチャート
- **サイドバー** — ボード一覧、イシュー数、最近のアイテム

### ユーザー体験
- **グローバル検索** — Cmd+K でクイック検索
- **キーボードショートカット**
  - `h` — ダッシュボードへ移動
  - `a` — 承認キューへ移動
  - `?` — ショートカット一覧の表示
  - `Cmd+S` — 保存
  - `Escape` — モーダル/ダイアログを閉じる
- **サイドバーの開閉** — 画面スペースの最適化
- **ダークモード** — テーマ切り替え
- **多言語対応** — 韓国語、英語、日本語、中国語をサポート

### 自動化 & 連携
- **Webhook** — Discord、Slack、Telegram に対応
  - イベント: `issue.created`、`issue.approved`、`issue.status_changed`、`issue.deleted`
- **JSONエクスポート** — 全イシューデータのエクスポート
- **右クリックコンテキストメニュー** — クイックアクションメニュー
- **イシューテンプレート** — ボードごとにlocalStorageへ保存

### MCPサーバー（AIエージェント連携）
10種類のツールでAIエージェントの自動化を実現:
1. `list_boards` — 全ボードの取得
2. `list_issues` — イシューの取得（boardId、statusでフィルタリング可能）
3. `get_issue` — イシュー詳細の取得（指示・基準・コメントを含む）
4. `create_issue` — 新規イシューの作成
5. `approve_issue` — Pending → Approved への承認
6. `claim_issue` — Approved → InProgress への遷移
7. `complete_issue` — InProgress → Done への完了（コメント任意）
8. `add_comment` — イシューへのコメント追加
9. `update_criteria` — 成功基準のチェック/アンチェック
10. `search_issues` — タイトルによるイシュー検索

## 技術スタック

### バックエンド
- **言語**: Go 1.26
- **ルーター**: chi/v5
- **データベース**: SQLite / PostgreSQL（選択可能）
- **マイグレーション**: goose/v3
- **ファイル埋め込み**: embed.FS（単一バイナリでデプロイ可能）

### フロントエンド
- **フレームワーク**: React 18
- **言語**: TypeScript
- **バンドラー**: Vite 6
- **CSS**: Tailwind CSS
- **状態管理**: React Query (TanStack)
- **ドラッグ**: @dnd-kit/core、@dnd-kit/sortable
- **アイコン**: lucide-react
- **チャート**: recharts
- **i18n**: react-i18next
- **Markdown**: react-markdown + rehype-sanitize

### MCP
- プロトコル: JSON-RPC 2.0 over stdio
- 対象: Claude Code、Claude Desktop

## はじめに

### 必要環境
- Go 1.26 以上
- Node.js 18 以上
- npm または yarn

### インストール & 起動

#### 1. リポジトリのクローンと初期化
```bash
git clone https://github.com/kjm99d/monkey-planner.git
cd monkey-planner
make init
```

#### 2. プロダクションビルド（単一バイナリ）
```bash
make build
./bin/monkey-planner
```

サーバーは `http://localhost:8080` で起動し、フロントエンドは内包されます。

#### 3. 開発モード（分離起動）

ターミナル1 — バックエンド:
```bash
make run-backend
```

ターミナル2 — フロントエンド（Vite dev server、:5173）:
```bash
make run-frontend
```

フロントエンドは `/api` へのリクエストを自動的に `:8080` へプロキシします。

### 環境変数

```bash
# サーバーアドレス（デフォルト: :8080）
export MP_ADDR=":8080"

# データベース接続文字列
# SQLite（デフォルト: sqlite://./data/monkey.db）
export MP_DSN="sqlite://./data/monkey.db"

# PostgreSQL の例
export MP_DSN="postgres://user:password@localhost:5432/monkey_planner"
```

## MCPサーバーの設定

### クイックセットアップ（自動更新）

1. [Releases](https://github.com/kjm99d/MonkeyPlanner/releases)から最新バイナリをダウンロード（例: `D:/mp/`）
2. 同じディレクトリに `update-and-run.sh` をダウンロード
3. `.mcp.json` または Claude Code 設定に追加:

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

ラッパースクリプトが起動時にGitHubの最新リリースを自動チェックし、バイナリを更新します。

### 手動セットアップ

自動更新なしで直接設定:

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

Claude Code（`.mcp.json`）と Claude Desktop（`~/.claude/claude_desktop_config.json`）の両方で同じ形式が使えます。変更後に再起動してください。

### MCPツールの使用例

```
AI: 全てのボードを一覧表示してください
→ list_boards() を呼び出し

AI: "認証" に関連するイシューを検索してください
→ search_issues(query="認証") を呼び出し

AI: 最初のPendingイシューを承認して進行中に変更し、完了させます
→ approve_issue() → claim_issue() → complete_issue() を順番に呼び出し
```

## エージェントのワークフロー

```
┌────────────────┐
│  人間がイシュー │  タイトル、本文、指示を入力
└────────┬───────┘
         │
         ↓
┌────────────────┐
│  Approveボタン │  Pending → Approved
└────────┬───────┘
         │
         ↓
┌────────────────────────────┐
│  AIエージェント（MCPクライアント）│  list_issues または search_issues
└────────┬───────────────────┘
         │
         ↓
┌────────────────────┐
│ claim_issue()      │  Approved → InProgress
└────────┬───────────┘
         │
         ↓
┌────────────────────┐
│ 作業を進行中...    │  add_comment()、update_criteria()
│                    │  （進捗報告 & 基準のチェック）
└────────┬───────────┘
         │
         ↓
┌────────────────────┐
│ complete_issue()   │  InProgress → Done
│ + 最終コメント     │
└────────┬───────────┘
         │
         ↓
┌────────────────┐
│  人間が確認    │  結果のレビューとフィードバック
└────────────────┘
```

## APIドキュメント

OpenAPI 3.0 仕様: [backend/docs/swagger.yaml](./backend/docs/swagger.yaml)

### 主要エンドポイント

#### Boards
```
GET    /api/boards                  # ボード一覧
POST   /api/boards                  # ボード作成
PATCH  /api/boards/{id}             # ボード更新
DELETE /api/boards/{id}             # ボード削除
```

#### Issues
```
GET    /api/issues                  # イシュー一覧（フィルタ: boardId、status、parentId）
POST   /api/issues                  # イシュー作成
GET    /api/issues/{id}             # イシュー詳細 + 子イシュー
PATCH  /api/issues/{id}             # イシュー更新（ステータス、属性、タイトルなど）
DELETE /api/issues/{id}             # イシュー削除
POST   /api/issues/{id}/approve     # イシュー承認（Pending → Approved）
```

#### Comments
```
GET    /api/issues/{issueId}/comments    # コメント一覧
POST   /api/issues/{issueId}/comments    # コメント追加
DELETE /api/comments/{commentId}         # コメント削除
```

#### Properties（カスタム属性）
```
GET    /api/boards/{boardId}/properties      # 属性定義の一覧
POST   /api/boards/{boardId}/properties      # 属性の作成
PATCH  /api/boards/{boardId}/properties/{propId}  # 属性の更新
DELETE /api/boards/{boardId}/properties/{propId}  # 属性の削除
```

#### Webhooks
```
GET    /api/boards/{boardId}/webhooks       # Webhook一覧
POST   /api/boards/{boardId}/webhooks       # Webhook作成
PATCH  /api/boards/{boardId}/webhooks/{whId}    # Webhook更新
DELETE /api/boards/{boardId}/webhooks/{whId}    # Webhook削除
```

#### Calendar
```
GET /api/calendar           # 月次統計（year、month が必須）
GET /api/calendar/day       # 日次イシュー一覧（date が必須）
```

詳細なスキーマは [backend/docs/swagger.yaml](./backend/docs/swagger.yaml) を参照してください。

## プロジェクト構成

```
monkey-planner/
├── backend/
│   ├── cmd/monkey-planner/
│   │   ├── main.go              # エントリーポイント（HTTPサーバー）
│   │   └── mcp.go               # MCPサーバー（JSON-RPC stdio）
│   ├── internal/
│   │   ├── domain/              # ドメインモデル（Issue、Boardなど）
│   │   ├── service/             # ビジネスロジック
│   │   ├── storage/             # データベース（SQLite/PostgreSQL）
│   │   ├── http/                # HTTPハンドラー & ルーター
│   │   └── migrations/          # goose マイグレーションファイル
│   ├── web/                     # フロントエンド埋め込み（embed.FS）
│   ├── docs/
│   │   └── swagger.yaml         # OpenAPI 3.0 仕様
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── components/          # 再利用可能なコンポーネント
│   │   ├── features/            # ページ & 機能別コンポーネント
│   │   │   ├── home/           # ダッシュボード
│   │   │   ├── board/          # ボード & カンバン
│   │   │   ├── issue/          # イシュー詳細
│   │   │   ├── calendar/       # カレンダー
│   │   │   └── approval/       # 承認キュー
│   │   ├── api/                 # APIフック & クライアント
│   │   ├── design/              # Tailwind トークン
│   │   ├── i18n/                # 多言語（en.json、ko.json、ja.json、zh.json）
│   │   ├── App.tsx              # ルーター
│   │   ├── index.css            # グローバルスタイル
│   │   └── main.tsx
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── tailwind.config.js
│
├── .mcp.json                    # Claude Code MCP 設定
├── Makefile                     # ビルド & 開発コマンド
├── .githooks/                   # Git フック
└── data/                        # SQLite データベース（デフォルト）
```

## テスト

### バックエンドテスト
```bash
make test-backend
```

### フロントエンドテスト
```bash
make test-frontend
```

### アクセシビリティテスト
```bash
make test-a11y
```

### 全テスト実行
```bash
make test
```

## 主要コマンド

```bash
# クローン後の初期セットアップ
make init

# プロダクションビルド
make build

# プロダクションサーバーの起動
./bin/monkey-planner

# 開発モード
make run-backend        # ターミナル1
make run-frontend       # ターミナル2

# クリーンアップ
make clean
```

## ステータス遷移ルール

```
Pending
  ↓ (approve endpoint)
Approved
  ↓ (PATCH status)
InProgress
  ↓ (PATCH status)
Done

Pending → Approved: POST /api/issues/{id}/approve （専用エンドポイント）
Approved ↔ InProgress ↔ Done: PATCH で自由に遷移可能
Pending: 他のステータスから戻ることはできない
Rejected: 却下ステータス（別途追跡）
```

## ライセンス

MIT

## コントリビューション

イシューやプルリクエストは歓迎します。

## お問い合わせ

プロジェクトに関するご質問やフィードバックは、GitHub Issues からお寄せください。
