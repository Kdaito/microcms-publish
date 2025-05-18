package md

import (
	"errors"
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
				Content: `<h2>これはテスト用の記事です。</h2>

<p>これはテスト用の記事です。</p>
`,
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
					Content: `<h2>これはテスト用の記事です。</h2>

<p>これはテスト用の記事です。</p>
`,
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
					Content: `<h2>これはテスト用の記事です。</h2>
<p>これはテスト用の記事です。</p>
`,
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
