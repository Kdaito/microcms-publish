package main

import (
	"log"
	// "os"
	// "os/exec"
	// "strings"
)

// func findModifiedMarkdownFiles(currentCommitHash string) ([]string, error) {
// 	cmd := exec.Command("git", "diff", "--name-only", currentCommitHash)
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Printf("Output: %s", output)

// 	var files []string
// 	// for _, file := range strings.Split(string(output), "\n") {
// 	// 	if strings.HasPrefix(file, "public/") && strings.HasSuffix(file, ".md") {
// 	// 		files = append(files, file)
// 	// 	}
// 	// }

// 	// for file := range string(output) {
// 	// 	log.Printf("Processing %s", file)
// 	// }
// 	return files, nil
// }

func main() {
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
