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

		expectedContent := `<h2>これはテスト用の記事です。</h2>

<p>これはテスト用の記事です。</p>
`

		// then
		assert.Equal(t, "テスト用の記事", item.Title)
		assert.Equal(t, "Test1,Test2", item.Tags)
		assert.Equal(t, "abcdefg12345", item.QiitaID)
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
		assert.Equal(t, "invalid front matter format", err.Error())
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
		assert.Equal(t, "invalid metadata format", err.Error())
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
file:[scanItem/invalidFrontMatter.md] scanning is skipped because: invalid front matter format
Scan: ../../mocks/scanItem/invalidMetadata.md
file:[scanItem/invalidMetadata.md] scanning is skipped because: invalid metadata format
Scan: ../../mocks/scanItem/withoutIdAndTitle.md
file:[scanItem/withoutIdAndTitle.md] scanning is skipped because: title or id is empty`
		assert.Equal(t, 1, len(items))
		actual := strings.TrimRight(buf.String(), "\n")
		assert.Equal(t, expectedLogs, actual)
	})
}
