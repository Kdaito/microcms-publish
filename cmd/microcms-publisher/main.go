package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Kdaito/microcms-publish/internal/cms"
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
	// 環境変数チェック
	var serviceId = os.Getenv(("SERVICE_ID"))
	if serviceId == "" {
		log.Fatal("SERVICE_ID is not set")
	}
	var apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY is not set")
	}
	var endpoint = os.Getenv("ENDPOINT")
	if endpoint == "" {
		log.Fatal("ENDPOINT is not set")
	}

	log.Printf("serviceId: %s", serviceId)
	// log.Printf("apiKey: %s", apiKey)
	// log.Printf("endpoint: %s", endpoint)

	// 差分のファイルを引数から取得する
	filesString := flag.String("f", "target files", "string array")
	workspace := flag.String("w", "workspace/path", "workspace path")
	flag.Parse()

	log.Printf("workspace: %s", *workspace)

	files := strings.Split(*filesString, ",")

	// ファイルから記事情報を取得する
	items := scanItems(&files, *workspace)

	if len(items) == 0 {
		log.Println("No items found.")
		return
	}

	httpClient := new(http.Client)

	// クライアントの初期化
	cmsClient := cms.NewClient(
		serviceId,
		apiKey,
		endpoint,
		httpClient,
	)

	// コンテキストの作成（タイムアウト付き）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 各記事をMicroCMSにアップロードする
	for _, item := range items {
		exists, id, err := cmsClient.CheckExists(ctx, item.QiitaID)
		if err != nil {
			log.Printf("Error checking existence: %v", err)
			continue
		}

		if exists {
			log.Printf("Content with ID %s already exists. Updating...", id)
			err = cmsClient.Update(ctx, id, item.Title, item.Tags, item.QiitaID, item.Content)
			if err != nil {
				log.Printf("Error updating content: %v", err)
			}
		} else {
			log.Println("Creating new content...")
			err = cmsClient.Create(ctx, item.Title, item.Tags, item.QiitaID, item.Content)
			if err != nil {
				log.Printf("Error creating content: %v", err)
			}
		}
	}
}
