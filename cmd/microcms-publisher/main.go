package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func scanItem(file string, workspace string) {
	log.Printf("Processing %s", file)
	log.Printf("Workspace: %s/%s", workspace, file)

	// ファイルの内容を取得する
	content, err := os.ReadFile(file)

	if err != nil {
		panic(err)
	}

	log.Printf("Content: %s", content)
}

func scanItems(files *[]string, workspace string) {
	items := make([]string, 0, len(*files))

	for _, file := range *files {
		scanItem(file, workspace)
		items = append(items, file)
	}
}

func main() {
	// 差分のファイルを引数から取得する
	filesString := flag.String("f", "target files", "string array")
	workspace := flag.String("w", "workspace/path", "workspace path")
	flag.Parse()

	log.Printf("workspace: %s", *workspace)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Println(dir)

	files := strings.Split(*filesString, ",")

	log.Printf("files: %s", files)

	scanItems(&files, *workspace)

	// currentCommitHash := os.Getenv("CURRENT_COMMIT_HASH")
	// _, err := findModifiedMarkdownFiles(currentCommitHash)
	// if err != nil {
	// 	panic(err)
	// }

	log.Printf("Hello, world from cmd !")

	// for _, file := range files {
	// 	log.Printf("Processing %s", file)
	// }
}
