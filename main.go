package main

import (
	"log"

	"autoCR/pkg/tool"
)

func main() {
	url := "https://github.com/owner/repo/pull/1" // Replace with a valid GitHub pull request URL
	err := tool.CodeReview(url)
	if err != nil {
		log.Fatalf("Code review failed: %v", err)
	}
	log.Println("Code review completed successfully")
}
