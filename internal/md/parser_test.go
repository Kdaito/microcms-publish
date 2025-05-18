package md

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewParser("test-workspace")
	if client == nil {
		t.Fatal("Expected client to be initialized, got nil")
	}
}

func TestParseFromQiitaItem(t *testing.T) {
	tests := []struct {
		name          string
		file          string
		expectedItem  *Item
		expectedError string
	}{
		{
			name: "正常系",
			file: "parseItem/success.md",
			expectedItem: &Item{
				Title:   "テスト用の記事",
				Tags:    "Test1,Test2",
				QiitaID: "abcdefg12345",
				Content: "<h2>これはテスト用の記事です。</h2>\n<p>これはテスト用の記事です。</p>\n",
			},
			expectedError: "",
		},
		{
			name:          "異常系_invalidFrontMatter",
			file:          "parseItem/invalidFrontMatter.md",
			expectedItem:  nil,
			expectedError: "invalid front matter format",
		},
		{
			name:          "異常系_invalidMetadata",
			file:          "parseItem/invalidMetadata.md",
			expectedItem:  nil,
			expectedError: "invalid metadata format",
		},
		{
			name:          "異常系_withoutIdAndTilte",
			file:          "parseItem/withoutIdAndTitle.md",
			expectedItem:  nil,
			expectedError: "title or id is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockFilePath := tt.file
			mockWorkspace := "../../mocks"

			// when
			parser := NewParser(mockWorkspace)
			item, err := parser.parseFromQiitaItem(mockFilePath)

			// then
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Nil(t, item)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedItem.Title, item.Title)
				assert.Equal(t, tt.expectedItem.Tags, item.Tags)
				assert.Equal(t, tt.expectedItem.QiitaID, item.QiitaID)
				assert.Equal(t, tt.expectedItem.Content, item.Content)
			}
		})
	}
}

// モック用のParser構造体
type MockParser struct {
	Parser
	mockParseFromQiitaItem func(file string) (*Item, error)
}

// モックのparseFromQiitaItemメソッドをオーバーライド
func (m *MockParser) parseFromQiitaItem(file string) (*Item, error) {
	if m.mockParseFromQiitaItem != nil {
		return m.mockParseFromQiitaItem(file)
	}
	return nil, errors.New("mock parse function not set")
}

func TestParseAllFromQiitaItems(t *testing.T) {
	tests := []struct {
		names         string
		files         []string
		expectedItems []*Item
	}{
		{
			names: "正常系",
			files: []string{
				"parseItem/success.md",
				"parseItem/invalidFrontMatter.md",
			},
			expectedItems: []*Item{
				{
					Title:   "テスト用の記事",
					Tags:    "Test1,Test2",
					QiitaID: "abcdefg12345",
					Content: "<h2>これはテスト用の記事です。</h2>\n<p>これはテスト用の記事です。</p>\n",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.names, func(t *testing.T) {
			// given
			mockWorkspace := "../../mocks"
			mockParser := &MockParser{
				Parser: Parser{
					workspace: mockWorkspace,
				},
			}
			mockParser.mockParseFromQiitaItem = func(file string) (*Item, error) {
				if strings.Contains(file, "invalidFrontMatter") {
					return nil, errors.New("invalid front matter format")
				}
				return &Item{
					Title:   "テスト用の記事",
					Tags:    "Test1,Test2",
					QiitaID: "abcdefg12345",
					Content: "<h2>これはテスト用の記事です。</h2>\n<p>これはテスト用の記事です。</p>\n",
				}, nil
			}

			// when
			items := mockParser.ParseAllFromQiitaItems(&tt.files)

			// then
			assert.Equal(t, len(tt.expectedItems), len(items))
			for i, expected := range tt.expectedItems {
				assert.Equal(t, expected.Title, items[i].Title)
				assert.Equal(t, expected.Tags, items[i].Tags)
				assert.Equal(t, expected.QiitaID, items[i].QiitaID)
				assert.Equal(t, expected.Content, items[i].Content)
			}
		})
	}
}

func TestParseHtml(t *testing.T) {
	tests := []struct {
		name           string
		targetFilePath string
		expected       string
	}{
		{
			name:           "正常系",
			targetFilePath: "../../mocks/parseHtml/success.md",
			expected:       "<h1>タイトル1です</h1>\n<p>ここは導入文です。この記事では、Markdownの基本文法について説明します。</p>\n<h2>タイトル2です</h2>\n<p>Markdownは<strong>シンプル</strong>で、_可読性_が高く、<code>コード</code>も簡単に書けます。</p>\n<h3>タイトル3です</h3>\n<p>以下にさまざまなMarkdownの構文を紹介します。</p>\n<hr>\n<h3>見出し</h3>\n<h1>見出し1</h1>\n<h2>見出し2</h2>\n<h3>見出し3</h3>\n<h4>見出し4</h4>\n<h5>見出し5</h5>\n<h6>見出し6</h6>\n<hr>\n<h3>リスト</h3>\n<ul>\n<li>箇条書き1\n<ul>\n<li>ネスト1\n<ul>\n<li>ネスト2</li>\n</ul>\n</li>\n</ul>\n</li>\n<li>箇条書き2</li>\n</ul>\n<ol>\n<li>番号付きリスト1</li>\n<li>番号付きリスト2\n<ol>\n<li>ネストされた番号付きリスト</li>\n</ol>\n</li>\n</ol>\n<hr>\n<h3>引用</h3>\n<blockquote>\n<p>これは引用です。<br>\n引用内で改行もできます。</p>\n</blockquote>\n<hr>\n<h3>コードブロック</h3>\n<h4>インラインコード</h4>\n<p>例えば、<code>console.log(&quot;Hello World&quot;)</code>のように書きます。</p>\n<h4>ブロックコード（シンタックスハイライト付き）</h4>\n<pre><code class=\"language-javascript\">function greet(name) {\n  console.log(`Hello, ${name}!`);\n}\ngreet(&quot;Markdown&quot;);\n</code></pre>\n<hr>\n<h3>テーブル</h3>\n<table>\n<thead>\n<tr>\n<th>名前</th>\n<th>年齢</th>\n<th>職業</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>山田太郎</td>\n<td>29</td>\n<td>エンジニア</td>\n</tr>\n<tr>\n<td>田中花子</td>\n<td>34</td>\n<td>デザイナー</td>\n</tr>\n</tbody>\n</table>\n<hr>\n<h3>リンクと画像</h3>\n<p><a href=\"https://www.google.com\">Google</a></p>\n<p><img src=\"https://images.dog.ceo/breeds/pembroke/n02113023_15998.jpg\" alt=\"ダミー画像\"></p>\n<hr>\n<h3>太字・斜体・打ち消し</h3>\n<ul>\n<li><strong>太字</strong></li>\n<li><em>斜体</em></li>\n<li>~~打ち消し~~</li>\n</ul>\n<hr>\n<h3>チェックリスト</h3>\n<ul>\n<li><input checked=\"\" disabled=\"\" type=\"checkbox\"> 記事構成を考える</li>\n<li><input disabled=\"\" type=\"checkbox\"> 実装する</li>\n<li><input disabled=\"\" type=\"checkbox\"> 公開する</li>\n</ul>\n<hr>\n<h3>改行の確認</h3>\n<p>この文の後には2スペースがあります。<br>\nなので改行されます。</p>\n<hr>\n<p>おわりに。この記事ではMarkdownの様々な構文を紹介しました。</p>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			content, err := os.ReadFile(tt.targetFilePath)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			parts := strings.SplitN(string(content), "---\n", 3)

			// when
			result := parseHtml(parts[2])
			assert.Equal(t, tt.expected, result)
		})
	}
}
