# MicroCMS Publish

MicroCMS Publish は、[qiita-cli](https://github.com/increments/qiita-cli) のカスタムアクションを拡張した GitHub Actions です。Qiita への記事投稿・更新と同時に、MicroCMS にも記事を反映できます。

## 機能

- `/public/xx.md` ファイルを変更・追加すると、Qiita に投稿し、その内容を MicroCMS にも反映
- Qiita の記事 ID (`qiitaId`) をキーとして MicroCMS に記事を作成・更新

## 事前準備

### 1. MicroCMS のセットアップ

MicroCMS に以下の API スキーマを持つ記事エンドポイントを作成してください。

| フィールド ID | 表示名（例）  | 種類               |
| ------------- | ------------- | ------------------ |
| title         | タイトル      | テキストフィールド |
| tags          | タグ          | テキストフィールド |
| qiitaId       | Qiita 記事 ID | テキストフィールド |
| content       | 記事本文      | リッチエディタ     |

### 2. リポジトリの構成

このアクションは、`qiita-cli` のリポジトリ構成に基づいて動作します。

「[Qiitaの記事をGitHubリポジトリで管理する方法](https://qiita.com/Qiita/items/32c79014509987541130)」を参考にリポジトリのセットアップを行なってください。

```
├─ .github/workflows/publish.yml
└─ public
   ├─ item001.md
   └─ itemxxx.md
```

環境変数には、以下を設定してください。

| 変数名                | 内容                                                                       |
| --------------------- | -------------------------------------------------------------------------- |
| `QIITA_TOKEN`         | Qiita の API トークン                                                      |
| `MICROCMS_API_KEY`    | MicroCMS の API キー                                                       |
| `MICROCMS_SERVICE_ID` | MicroCMS のサービス ID <br/>（例: `https://Kdaito.microcms.io` の場合 `Kdaito`） |

## GitHub Actions の設定

`.github/workflows/publish.yml` を作成し、以下のジョブを追加してください。

```yaml
jobs:
  publish_articles:
    runs-on: ubuntu-latest
    timeout-minutes: 5
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: Kdaito/microcms-publish/actions/publish@main
        with:
          api-key: ${{ secrets.MICROCMS_API_KEY }}
          service-id: ${{ secrets.MICROCMS_SERVICE_ID }}
          qiita-token: ${{ secrets.QIITA_TOKEN }}
          endpoint: "items"
```

`endpoint` には、MicroCMS で作成したエンドポイントの ID を指定してください。

## 投稿方法

[qiita-cli](https://github.com/increments/qiita-cli) を使用して GitHub で Qiita の記事を管理する場合と同様の運用が可能です。

参考: [Qiita の記事を GitHub リポジトリで管理する方法](https://qiita.com/Qiita/items/32c79014509987541130)

### 記事の例

以下の `.md` ファイルをリポジトリに追加すると、Qiita と MicroCMS に投稿されます。

```md
---
title: サンプル記事タイトル
tags:
  - Java
  - TypeScript
  - 型
private: false
updated_at: "2025-02-10T23:57:01+09:00"
id: 12345abcde
organization_url_name: null
slide: false
ignorePublish: false
---

## タイトル

内容
```

MicroCMS では、以下のようにデータが登録されます。

| フィールド ID | 内容                           |
| ------------- | ------------------------------ |
| title         | サンプル記事タイトル           |
| tags          | Java,TypeScript,型             |
| qiitaId       | 12345abcde                     |
| content       | `<h2>タイトル</h2><p>内容</p>` |

新規投稿時は、`id` を `null` に設定してください。Qiita CLI のカスタムアクションが ID を付与した後、MicroCMS に反映されます。
