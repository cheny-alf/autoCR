package github

import (
	"context"
	"fmt"
	"testing"
)

func mockClient() *GithubClient {
	return NewGithubClient("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").SetUserInfo("cheny-alf", "autoCR")
}

func TestGetCommitFileContent(t *testing.T) {
	client := mockClient()
	content, err := client.GetCommitFileContent(context.Background(), "go.mod", "fc22399")
	if err != nil {
		fmt.Println(err)
	}
	t.Log(content)
}

func TestGetPendingReviewCommits(t *testing.T) {
	client := mockClient()
	fmt.Print(client)
	list, err := client.GetPendingReviewCommits(context.Background(), 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(list)
}

func TestGetPullRequestFiles(t *testing.T) {
	client := mockClient()
	err := client.GetPullRequestFiles(context.Background(), 1)
	if err != nil {
		t.Error(err)
	}
	t.Log()
}
