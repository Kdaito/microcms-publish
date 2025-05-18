package md

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/russross/blackfriday/v2"
)

type QiitaItemMetadata struct {
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
	Id    string   `yaml:"id"`
}

type Item struct {
	Title   string `json:"title"`
	Tags    string `json:"tags"`
	QiitaID string `json:"qiitaId"`
	Content string `json:"content"`
}

type Parser struct {
	workspace string
}

func NewParser(workspace string) *Parser {
	return &Parser{
		workspace: workspace,
	}
}

func (s *Parser) parseFromQiitaItem(file string) (*Item, error) {
	filePath := fmt.Sprintf("%s/%s", s.workspace, file)

	log.Printf("Parse: %s", filePath)

	// ファイルの内容を取得する
	content, err := os.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	parts := strings.SplitN(string(content), "---\n", 3)
	if len(parts) < 3 {
		return nil, errors.New("invalid front matter format")
	}

	var qiitaItemMetadata QiitaItemMetadata
	if err := yaml.Unmarshal([]byte(parts[1]), &qiitaItemMetadata); err != nil {
		return nil, errors.New("invalid metadata format")
	}

	if qiitaItemMetadata.Title == "" || qiitaItemMetadata.Id == "" {
		return nil, errors.New("title or id is empty")
	}

	htmlContent := string(blackfriday.Run([]byte(parts[2])))
	item := &Item{
		Title:   qiitaItemMetadata.Title,
		Tags:    strings.Join(qiitaItemMetadata.Tags, ","),
		QiitaID: qiitaItemMetadata.Id,
		Content: htmlContent,
	}

	return item, nil
}

func (s *Parser) ParseAllFromQiitaItems(files *[]string) []*Item {
	items := make([]*Item, 0, len(*files))

	for _, file := range *files {
		item, err := s.parseFromQiitaItem(file)
		if err != nil {
			log.Printf("file:[%s] parsing is skipped because: %s", file, err)
			continue
		}
		items = append(items, item)
	}

	return items
}
