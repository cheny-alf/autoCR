package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v69/github"
)

const (
	X_GitHub_Api_Version = "2022-11-28"
	Accept               = "application/vnd.github+json"
	token                = ""
)

type GithubClient struct {
	Client      *github.Client
	Owner, Repo string
}

func NewGithubClient(token string) *GithubClient {
	return &GithubClient{
		Client: github.NewClient(nil).WithAuthToken(token),
	}
}

func (g *GithubClient) SetUserInfo(owner, repo string) *GithubClient {
	g.Owner = owner
	g.Repo = repo
	return g
}

// 根据pr获取相关地址
func getPRDetail(url string) {

}

func checkRepo() {

}

// 获取指定commit中 文件内容
func (g *GithubClient) GetCommitFileContent(ctx context.Context, filePath, sha string) (string, error) {
	fileContent, _, _, err := g.Client.Repositories.GetContents(ctx, g.Owner, g.Repo, filePath, &github.RepositoryContentGetOptions{Ref: sha})
	if err != nil {
		return "", err
	}
	content, err := fileContent.GetContent()
	// 补充行数
	var contentWithLine []string
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		contentWithLine = append(contentWithLine, fmt.Sprintf("%d %s", i+1, line))
	}

	return strings.Join(contentWithLine, "\n"), err
}

func (g *GithubClient) GetPendingReviewCommits(ctx context.Context, prId int) ([]*github.RepositoryCommit, error) {
	// 获取commit
	commits, _, err := g.Client.PullRequests.ListCommits(ctx, g.Owner, g.Repo, prId, nil)
	if err != nil {
		return nil, err
	}

	var pendingReviewCommits []*github.RepositoryCommit
	for _, commit := range commits {
		if strings.HasPrefix(*commit.Commit.Message, "Merge pull request") {
			continue
		}
		commitInfo, _, err := g.Client.Repositories.GetCommit(ctx, g.Owner, g.Repo, *commit.SHA, nil)
		if err != nil {
			return nil, err
		}
		fmt.Println(commitInfo.Files)
		fmt.Println("(-------------)")

		pendingReviewCommits = append(pendingReviewCommits, commitInfo)
	}
	// fmt.Println(pendingReviewCommits)
	return pendingReviewCommits, nil
}

// 创建评论
func CreateCommitComment() {

}
