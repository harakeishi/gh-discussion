# GitHub CLI拡張 `gh-discussion` 実装計画

## 概要
GitHub CLI拡張として、GitHubディスカッションの検索と内容取得機能を提供する`gh-discussion`を実装する。

## 機能要件

### サブコマンド構造（gh issue/pr に準拠）

#### GENERAL COMMANDS
- **list**: ディスカッション一覧表示（`gh issue list`に相当）
- **view**: 特定ディスカッションの詳細表示（`gh issue view`に相当）
- **create**: 新しいディスカッション作成（`gh issue create`に相当）

#### TARGETED COMMANDS（将来実装）
- **comment**: ディスカッションにコメント追加
- **edit**: ディスカッション編集
- **close**: ディスカッション終了
- **reopen**: ディスカッション再開
- **lock/unlock**: ディスカッション会話のロック/ロック解除

### 1. ディスカッション一覧機能（list）
- **基本フィルタ**: `-a, --author`, `-S, --search`, `-l, --label`
- **ディスカッション固有フィルタ**: `--category`, `--answered`, `--unanswered`
- **共通オプション**: `-L, --limit`, `--json`, `-w, --web`
- **出力形式**: デフォルトはテーブル形式、`--json`で JSON 出力

### 2. ディスカッション詳細表示機能（view）
- **識別子**: ディスカッション番号またはURL
- **詳細表示**: `-c, --comments`でコメント表示
- **出力形式**: `--json`, `-w, --web`オプション対応

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
│   ├── list.go            # ディスカッション一覧コマンド実装
│   ├── view.go            # ディスカッション詳細表示コマンド実装
│   ├── create.go          # ディスカッション作成コマンド実装
│   └── comment.go         # コメント追加コマンド実装（将来実装）
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

### Phase 4: list機能実装
1. ディスカッション一覧取得の実装
2. フィルタ条件のバリデーション
3. 検索結果の取得と処理

### Phase 5: view機能実装
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

### list機能の GraphQL クエリ例

```graphql
query ListDiscussions($owner: String!, $repo: String!, $first: Int!, $after: String, $orderBy: DiscussionOrder, $filterBy: DiscussionOrderField) {
  repository(owner: $owner, name: $repo) {
    discussions(first: $first, after: $after, orderBy: $orderBy, filterBy: $filterBy) {
      pageInfo {
        hasNextPage
        endCursor
      }
      nodes {
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
        url
        answerChosenAt
        isAnswered
        upvoteCount
        comments(first: 0) {
          totalCount
        }
        labels(first: 10) {
          nodes {
            name
            color
          }
        }
      }
    }
  }
}
```

### 検索機能の GraphQL クエリ例（search使用）

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

### view機能の GraphQL クエリ例

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

### listコマンド（gh issue list準拠）
```bash
# 基本的な一覧表示
gh discussion list

# リポジトリ指定
gh discussion list -R owner/repo

# 作成者でフィルタ
gh discussion list -a username

# 検索クエリでフィルタ
gh discussion list -S "API documentation"

# カテゴリでフィルタ
gh discussion list --category "General"

# 回答済み/未回答でフィルタ
gh discussion list --answered
gh discussion list --unanswered

# 表示件数制限
gh discussion list -L 50

# JSON出力（特定フィールドのみ）
gh discussion list --json "number,title,author,category,isAnswered,createdAt,updatedAt,url"

# JSON出力（全フィールド）
gh discussion list --json

# JQフィルタ適用
gh discussion list --jq '.[] | select(.isAnswered == false)'

# テンプレート出力
gh discussion list --template '{{range .}}{{.number}} {{.title}} by {{.author.login}}{{"\n"}}{{end}}'

# ブラウザで開く
gh discussion list -w

# 複合条件
gh discussion list -a username --category "General" -S "bug" --answered
```

### viewコマンド（gh issue view準拠）
```bash
# 番号で表示
gh discussion view 123

# URLで表示
gh discussion view https://github.com/owner/repo/discussions/123

# リポジトリ指定
gh discussion view 123 -R owner/repo

# コメントも表示
gh discussion view 123 -c

# JSON出力（特定フィールドのみ）
gh discussion view 123 --json "title,body,author,category,comments,isAnswered,createdAt"

# JSON出力（全フィールド）
gh discussion view 123 --json

# JQフィルタ適用
gh discussion view 123 --jq '.comments[] | select(.isAnswer == true)'

# テンプレート出力
gh discussion view 123 --template '{{.title}} by {{.author.login}} in {{.category.name}}'

# ブラウザで開く
gh discussion view 123 -w
```

### createコマンド（将来実装、gh issue create準拠）
```bash
# インタラクティブ作成
gh discussion create

# タイトルと本文指定
gh discussion create --title "Discussion Title" --body "Discussion body"

# カテゴリ指定
gh discussion create --category "General"
```

## JSON FIELDS

### listコマンド用JSON出力フィールド
```
activeLockReason, answer, answerChosenAt, answerChosenBy, author, authorAssociation,
body, bodyHTML, bodyText, category, comments, createdAt, createdViaEmail, databaseId,
editor, id, includesCreatedEdit, isAnswered, lastEditedAt, locked, number, 
publishedAt, reactionGroups, reactions, repository, resourcePath, title, updatedAt,
url, userContentEdits, viewerCanDelete, viewerCanReact, viewerCanSubscribe,
viewerCanUpdate, viewerDidAuthor, viewerSubscription
```

### viewコマンド用JSON出力フィールド
```
activeLockReason, answer, answerChosenAt, answerChosenBy, author, authorAssociation,
body, bodyHTML, bodyText, category, comments, createdAt, createdViaEmail, databaseId,
editor, id, includesCreatedEdit, isAnswered, lastEditedAt, locked, number, 
publishedAt, reactionGroups, reactions, repository, resourcePath, title, updatedAt,
url, userContentEdits, viewerCanDelete, viewerCanReact, viewerCanSubscribe,
viewerCanUpdate, viewerDidAuthor, viewerSubscription
```

### ネストされたオブジェクトのフィールド

#### author, answerChosenBy, editor
```
avatarUrl, login, url, id, name, email
```

#### category
```
id, name, description, emoji, emojiHTML, isAnswerable, createdAt, updatedAt
```

#### comments
```
author, authorAssociation, body, bodyHTML, bodyText, createdAt, id, isAnswer,
isMinimized, minimizedReason, publishedAt, reactionGroups, replies, replyTo,
updatedAt, url, viewerCanMarkAsAnswer, viewerCanUnmarkAsAnswer
```

#### repository
```
id, name, nameWithOwner, owner, url, description
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
	github.com/cli/go-gh v1.2.1           // GitHub CLI 統合
	github.com/spf13/cobra v1.8.0          // CLI フレームワーク  
	github.com/cli/safeexec v1.0.1         // 安全なコマンド実行
	github.com/MakeNowJust/heredoc v1.0.0  // ヒアドキュメント
	github.com/briandowns/spinner v1.23.0  // スピナー表示
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
