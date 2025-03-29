package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanItem(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		// given
		mockFilePath := "scanItem/success.md"
		mockWorkspace := "../../mocks"

		// when
		item, err := scanItem(mockFilePath, mockWorkspace)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expectedContent := `<h2>そもそもシングルトンとは？</h2>

<p>インスタンスが必ず一つであることを保証するようなデザインパターンの一つです。</p>

<h2>Node.js におけるシングルトンの実装方法</h2>

<p>node.js では、インスタンスをエクスポートするだけでシングルトンを実現できるみたいです。</p>

<pre><code class="language-javascript">class SignletonClass {
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
</code></pre>

<p>Node.js のドキュメントでは、<code>Caching</code>にて以下のような説明がなされています。</p>

<blockquote>
<p>Modules are cached after the first time they are loaded. This means (among other things) that every call to require(&lsquo;foo&rsquo;) will get exactly the same object returned, if it would resolve to the same file.</p>
</blockquote>

<p>翻訳すると、</p>

<blockquote>
<p>モジュールは、最初にロードされた後にキャッシュされます。これは、(とりわけ) require(&lsquo;foo&rsquo;) を呼び出すたびに、同じファイル名の場合、まったく同じオブジェクトが返されることを意味します。</p>
</blockquote>

<p>なので一度インスタンスにしてしまってエクスポートしてしまえば、インポートするたびに新たなインスタンスは作成されず、最初に作成されたインスタンスが使いまわされるっぽいです。</p>

<p>引用: <a href="https://nodejs.org/api/modules.html#caching">https://nodejs.org/api/modules.html#caching</a></p>

<p>使いまわしたくない場合は、クラスをエクスポートして、インポート先でインスタンス化するようにしましょう。</p>

<h2>まとめ</h2>

<p>Node.js 環境では、インスタンスをエクスポートして使いまわせばシングルトンのように扱うことができる。</p>

<p>まちがってたらごめんなさい。</p>
`

		// then
		assert.Equal(t, "Node.jsにおけるシングルトン", item.Title)
		assert.Equal(t, "JavaScript,Node.js", item.Tags)
		assert.Equal(t, "d236248eb2d41173a96f", item.QiitaID)
		assert.Equal(t, expectedContent, item.Content)
	})

	t.Run("異常系_invalidFrontMatter", func(t *testing.T) {
		// given
		mockFilePath := "scanItem/invalidFrontMatter.md"
		mockWorkspace := "../../mocks"

		// when
		item, err := scanItem(mockFilePath, mockWorkspace)

		// then
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Equal(t, "invalid front matter format in ../../mocks/scanItem/invalidFrontMatter.md", err.Error())
	})

	t.Run("異常系_invalidMetadata", func(t *testing.T) {
		// given
		mockFilePath := "scanItem/invalidMetadata.md"
		mockWorkspace := "../../mocks"

		// when
		item, err := scanItem(mockFilePath, mockWorkspace)

		// then
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Equal(t, "invalid metadata format error unmarshaling JSON: json: cannot unmarshal string into Go struct field Metadata.Tags of type []string", err.Error())
	})
	t.Run("異常系_withoutIdAndTilte", func(t *testing.T) {
		// given
		mockFilePath := "scanItem/withoutIdAndTitle.md"
		mockWorkspace := "../../mocks"

		// when
		item, err := scanItem(mockFilePath, mockWorkspace)

		// then
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Equal(t, "title or id is empty", err.Error())
	})
}

func TestScanItems(t *testing.T) {
	// ログの出力先をバッファに変更
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// デフォルトだと日付が出力されてしまうので、フラグに0を設定する
	defaultFlags := log.Flags()
	log.SetFlags(0)

	// テスト終了時、変更した内容を戻す
	defer func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(defaultFlags)
	}()

	t.Run("正常系", func(t *testing.T) {
		defer func() {
			buf.Reset()
		}()

		// given
		mockFiles := []string{
			"scanItem/success.md",
			"scanItem/invalidFrontMatter.md",
			"scanItem/invalidMetadata.md",
			"scanItem/withoutIdAndTitle.md",
		}
		mockWorkspace := "../../mocks"

		// when
		items := scanItems(&mockFiles, mockWorkspace)

		// then
		expectedLogs := `Scan: ../../mocks/scanItem/success.md
Scan: ../../mocks/scanItem/invalidFrontMatter.md
file:[scanItem/invalidFrontMatter.md] scanning is skipped because: invalid front matter format in ../../mocks/scanItem/invalidFrontMatter.md
Scan: ../../mocks/scanItem/invalidMetadata.md
file:[scanItem/invalidMetadata.md] scanning is skipped because: invalid metadata format error unmarshaling JSON: json: cannot unmarshal string into Go struct field Metadata.Tags of type []string
Scan: ../../mocks/scanItem/withoutIdAndTitle.md
file:[scanItem/withoutIdAndTitle.md] scanning is skipped because: title or id is empty`
		assert.Equal(t, 1, len(items))
		actual := strings.TrimRight(buf.String(), "\n")
		assert.Equal(t, expectedLogs, actual)
	})
}
