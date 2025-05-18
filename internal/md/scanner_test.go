package md

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewScanner("test-workspace")
	if client == nil {
		t.Fatal("Expected client to be initialized, got nil")
	}
}

func TestScanFromQiitaItem(t *testing.T) {
	tests := []struct {
		name          string
		file          string
		expectedItem  *Item
		expectedError string
	}{
		{
			name: "正常系",
			file: "scanItem/success.md",
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
			file:          "scanItem/invalidFrontMatter.md",
			expectedItem:  nil,
			expectedError: "invalid front matter format",
		},
		{
			name:          "異常系_invalidMetadata",
			file:          "scanItem/invalidMetadata.md",
			expectedItem:  nil,
			expectedError: "invalid metadata format",
		},
		{
			name:          "異常系_withoutIdAndTilte",
			file:          "scanItem/withoutIdAndTitle.md",
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
			scanner := NewScanner(mockWorkspace)
			item, err := scanner.scanFromQiitaItem(mockFilePath)

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

// モック用のScanner構造体
type MockScanner struct {
	Scanner
	mockScanFromQiitaItem func(file string) (*Item, error)
}

// モックのscanFromQiitaItemメソッドをオーバーライド
func (m *MockScanner) scanFromQiitaItem(file string) (*Item, error) {
	if m.mockScanFromQiitaItem != nil {
		return m.mockScanFromQiitaItem(file)
	}
	return nil, errors.New("mock scan function not set")
}

func TestScanAllFromQiitaItems(t *testing.T) {
	tests := []struct {
		names         string
		files         []string
		expectedItems []*Item
	}{
		{
			names: "正常系",
			files: []string{
				"scanItem/success.md",
				"scanItem/invalidFrontMatter.md",
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
			mockScanner := &MockScanner{
				Scanner: Scanner{
					workspace: mockWorkspace,
				},
			}
			mockScanner.mockScanFromQiitaItem = func(file string) (*Item, error) {
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
			items := mockScanner.ScanAllFromQiitaItems(&tt.files)

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
