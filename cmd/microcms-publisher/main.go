package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/russross/blackfriday/v2"
)

type Metadata struct {
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

func scanItem(file string, workspace string) (*Item, error) {
	filePath := fmt.Sprintf("%s/%s", workspace, file)

	log.Printf("Scan: %s", filePath)

	// ファイルの内容を取得する
	content, err := os.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	parts := strings.SplitN(string(content), "---\n", 3)
	if len(parts) < 3 {
		return nil, errors.New("invalid front matter format")
	}

	var metadata Metadata
	if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
		return nil, errors.New("invalid metadata format")
	}

	if metadata.Title == "" || metadata.Id == "" {
		return nil, errors.New("title or id is empty")
	}

	htmlContent := string(blackfriday.Run([]byte(parts[2])))
	item := &Item{
		Title:   metadata.Title,
		Tags:    strings.Join(metadata.Tags, ","),
		QiitaID: metadata.Id,
		Content: htmlContent,
	}

	return item, nil
}



func scanItems(files *[]string, workspace string) []*Item {
	items := make([]*Item, 0, len(*files))

	for _, file := range *files {
		itme, err := scanItem(file, workspace)
		if err != nil {
			log.Printf("file:[%s] scanning is skipped because: %s", file, err)
			continue
		}
		items = append(items, itme)
	}

	return items
}

func main() {
	// 差分のファイルを引数から取得する
	filesString := flag.String("f", "target files", "string array")
	workspace := flag.String("w", "workspace/path", "workspace path")
	flag.Parse()

	log.Printf("workspace: %s", *workspace)

	files := strings.Split(*filesString, ",")

	items := scanItems(&files, *workspace)
}
