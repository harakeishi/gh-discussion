# GitHub CLI拡張 `gh-discussion` 実装計画

## 概要
GitHub CLI拡張として、GitHubディスカッションの検索と内容取得機能を提供する`gh-discussion`を実装する。

## 機能要件

### 1. ディスカッション検索機能
- **リポジトリ指定（必須）**: `-r, --repo` オプション
- **期間指定（オプション）**: `--created-after`, `--created-before`, `--updated-after`, `--updated-before`
- **ユーザー指定（オプション）**: `--author`, `--commenter`
- **キーワード検索（オプション）**: `--query`
- **カテゴリ指定（オプション）**: `--category`
- **出力形式**: デフォルトはテーブル形式、`--json`で JSON 出力

### 2. ディスカッション内容取得機能
- **ディスカッション ID 指定**: `--discussion-id`
- **詳細情報取得**: タイトル、本文、コメント、作成者情報など
- **出力形式**: デフォルトは読みやすい形式、`--json`で JSON 出力

## 技術スタック

- **実装言語**: Go（クロスプラットフォーム対応、高性能）
- **API**: GitHub GraphQL API v4
- **認証**: GitHub CLI の認証情報を利用
- **依存関係管理**: Go modules

## プロジェクト構造

```
gh-discussion/
├── main.go                 # エントリーポイント
├── cmd/
│   ├── search.go          # 検索コマンド実装
│   └── get.go             # 取得コマンド実装
├── pkg/
│   ├── client/
│   │   └── github.go      # GraphQL クライアント
│   ├── models/
│   │   └── discussion.go  # データモデル
│   └── formatter/
│       └── output.go      # 出力フォーマッター
├── go.mod
├── go.sum
├── README.md
└── .github/
    └── workflows/
        └── release.yml     # リリース自動化
```

## 実装ステップ

### Phase 1: プロジェクト初期化
1. Go プロジェクトの初期化
2. 必要な依存関係の追加
3. GitHub CLI 拡張としての設定

### Phase 2: GraphQL クライアント実装
1. GitHub GraphQL API クライアントの実装
2. 認証処理の実装（GitHub CLI の認証情報利用）
3. エラーハンドリングの実装

### Phase 3: データモデル定義
1. Discussion 構造体の定義
2. GraphQL レスポンスのマッピング
3. 検索条件の構造体定義

### Phase 4: 検索機能実装
1. 検索クエリビルダーの実装
2. 検索条件のバリデーション
3. 検索結果の取得と処理

### Phase 5: 内容取得機能実装
1. 特定ディスカッション取得の実装
2. コメント情報の取得
3. 詳細情報の整形

### Phase 6: CLI インターフェース実装
1. cobra ライブラリを使用したCLI構築
2. コマンドライン引数の解析
3. ヘルプメッセージの実装

### Phase 7: 出力フォーマッター実装
1. テーブル形式出力の実装
2. JSON 形式出力の実装
3. カラー対応とスタイリング

### Phase 8: テストとドキュメント
1. ユニットテストの作成
2. 統合テストの実装
3. README とドキュメントの作成

## API 仕様

### 検索機能の GraphQL クエリ例

```graphql
query SearchDiscussions($query: String!, $first: Int!, $after: String) {
  search(type: DISCUSSION, query: $query, first: $first, after: $after) {
    pageInfo {
      hasNextPage
      endCursor
    }
    nodes {
      ... on Discussion {
        id
        number
        title
        bodyText
        createdAt
        updatedAt
        author {
          login
          url
        }
        category {
          name
        }
        repository {
          nameWithOwner
        }
        url
        answerChosenAt
        isAnswered
        upvoteCount
        comments(first: 0) {
          totalCount
        }
      }
    }
  }
}
```

### 内容取得の GraphQL クエリ例

```graphql
query GetDiscussion($owner: String!, $repo: String!, $number: Int!) {
  repository(owner: $owner, name: $repo) {
    discussion(number: $number) {
      id
      number
      title
      body
      bodyText
      createdAt
      updatedAt
      author {
        login
        url
        avatarUrl
      }
      category {
        name
        description
      }
      labels(first: 10) {
        nodes {
          name
          color
        }
      }
      url
      locked
      answerChosenAt
      isAnswered
      upvoteCount
      comments(first: 100) {
        nodes {
          id
          body
          createdAt
          author {
            login
            url
          }
          upvoteCount
          isAnswer
        }
        pageInfo {
          hasNextPage
          endCursor
        }
      }
    }
  }
}
```

## コマンド仕様

### 検索コマンド
```bash
# 基本的な検索
gh discussion search -r owner/repo

# 期間指定検索
gh discussion search -r owner/repo --created-after 2024-01-01 --created-before 2024-12-31

# ユーザー指定検索
gh discussion search -r owner/repo --author username

# キーワード検索
gh discussion search -r owner/repo --query "API documentation"

# 複合条件検索
gh discussion search -r owner/repo --author username --category "General" --query "bug"

# JSON 出力
gh discussion search -r owner/repo --json
```

### 内容取得コマンド
```bash
# 特定ディスカッションの取得
gh discussion get -r owner/repo --discussion-id 123

# JSON 出力
gh discussion get -r owner/repo --discussion-id 123 --json
```

## エラーハンドリング

### 想定されるエラー
1. **認証エラー**: GitHub CLI が未ログインまたは権限不足
2. **API エラー**: レート制限、ネットワークエラー
3. **入力エラー**: 不正なリポジトリ名、存在しないディスカッション ID
4. **権限エラー**: プライベートリポジトリへのアクセス権限不足

### エラー対応
- 分かりやすいエラーメッセージの表示
- 適切な終了コードの設定
- デバッグ情報の出力（`--debug` フラグ）

## 依存関係

```go
module github.com/harakeishi/gh-discussion

go 1.21

require (
	github.com/cli/go-gh v1.2.1
	github.com/spf13/cobra v1.8.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/fatih/color v1.16.0
)
```

## リリース戦略

1. **バージョニング**: Semantic Versioning (SemVer) を採用
2. **リリース自動化**: GitHub Actions で自動ビルド・リリース
3. **マルチプラットフォーム**: Windows, macOS, Linux 対応
4. **インストール**: `gh extension install harakeishi/gh-discussion`

## セキュリティ考慮事項

1. GitHub CLI の認証情報の安全な利用
2. 入力値のサニタイゼーション
3. GraphQL インジェクション対策
4. 機密情報のログ出力防止

## パフォーマンス考慮事項

1. GraphQL クエリの最適化
2. ページネーション対応
3. 結果のキャッシュ機能（将来実装）
4. 並行処理による高速化

## 今後の拡張予定

1. **インタラクティブモード**: TUI での操作
2. **設定ファイル**: デフォルト値の設定
3. **出力テンプレート**: カスタム出力形式
4. **webhooks 連携**: ディスカッション更新の通知 
