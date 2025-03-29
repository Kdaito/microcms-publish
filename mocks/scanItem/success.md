---
title: Node.jsにおけるシングルトン
tags:
  - JavaScript
  - Node.js
private: false
updated_at: '2025-03-23T20:50:41+09:00'
id: d236248eb2d41173a96f
organization_url_name: null
slide: false
ignorePublish: false
---

## そもそもシングルトンとは？

インスタンスが必ず一つであることを保証するようなデザインパターンの一つです。

## Node.js におけるシングルトンの実装方法

node.js では、インスタンスをエクスポートするだけでシングルトンを実現できるみたいです。

```javascript
class SignletonClass {
  item: number;

  constructor() {
    this.item = 0;
  }

  get() {
    return this.item;
  }
}

const singletonClass = new SignletonClass();

// クラスではなく、インスタンスをエクスポート
export default singletonClass;
```

Node.js のドキュメントでは、`Caching`にて以下のような説明がなされています。

> Modules are cached after the first time they are loaded. This means (among other things) that every call to require('foo') will get exactly the same object returned, if it would resolve to the same file.

翻訳すると、

> モジュールは、最初にロードされた後にキャッシュされます。これは、(とりわけ) require('foo') を呼び出すたびに、同じファイル名の場合、まったく同じオブジェクトが返されることを意味します。

なので一度インスタンスにしてしまってエクスポートしてしまえば、インポートするたびに新たなインスタンスは作成されず、最初に作成されたインスタンスが使いまわされるっぽいです。

引用: https://nodejs.org/api/modules.html#caching

使いまわしたくない場合は、クラスをエクスポートして、インポート先でインスタンス化するようにしましょう。

## まとめ

Node.js 環境では、インスタンスをエクスポートして使いまわせばシングルトンのように扱うことができる。

まちがってたらごめんなさい。