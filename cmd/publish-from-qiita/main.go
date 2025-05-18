package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Kdaito/microcms-publish/internal/cms"
	"github.com/Kdaito/microcms-publish/internal/md"
)

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

	// 差分のファイルを引数から取得する
	filesString := flag.String("f", "target files", "string array")
	workspace := flag.String("w", "workspace/path", "workspace path")
	flag.Parse()

	log.Printf("workspace: %s", *workspace)

	files := strings.Split(*filesString, ",")

	// ファイルから記事情報を取得する
	parser := md.NewParser(*workspace)
	items := parser.ParseAllFromQiitaItems(&files)

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

	successItems := make([]string, 0, len(items))

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
			successItems = append(successItems, item.QiitaID)
		} else {
			log.Println("Creating new content...")
			err = cmsClient.Create(ctx, item.Title, item.Tags, item.QiitaID, item.Content)
			if err != nil {
				log.Printf("Error creating content: %v", err)
			}
			successItems = append(successItems, item.QiitaID)
		}
	}

	log.Println("Successfully processed items:")
	for _, id := range successItems {
		log.Println(id)
	}
	log.Println("All items processed.")
	log.Println("Publishing completed.")
}
