package tool

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	autoCRgithub "autoCR/pkg/tool/github"

	"github.com/google/go-github/v69/github"
)

type Reviewer interface {
	GetCommitFileContent(ctx context.Context, filePath, sha string) (string, error)
	GetPendingReviewCommits() ([]*github.RepositoryCommit, error)
	GetPullRequestReviews(ctx context.Context) ([]*github.PullRequestReview, error)
	SubmitCodeReview(ctx context.Context, commit github.RepositoryCommit, comments []*github.DraftReviewComment) error
}

// 执行入口
func CodeReview(url string) error {
	// 解析URL以提取所有者、仓库和PR编号
	repoURLRegex := regexp.MustCompile(`https://github\.com/([^/]+)/([^/]+)/pull/(\d+)`)
	matches := repoURLRegex.FindStringSubmatch(url)
	if len(matches) != 4 {
		return fmt.Errorf("invalid GitHub pull request URL")
	}
	owner, repo, prIDStr := matches[1], matches[2], matches[3]
	prID, err := strconv.Atoi(prIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse PR ID: %w", err)
	}

	// 初始化GitHub客户端
	client := autoCRgithub.NewGithubClient().SetUserInfo(owner, repo)

	// 获取待审核的提交
	pendingCommits, err := client.GetPendingReviewCommits(context.Background(), prID)
	if err != nil {
		return fmt.Errorf("failed to get pending review commits: %w", err)
	}

	// 遍历每个提交并获取文件内容
	var comments []*github.DraftReviewComment
	for _, commit := range pendingCommits {
		files, _, err := client.Client.PullRequests.ListFiles(context.Background(), owner, repo, prID, nil)
		if err != nil {
			return fmt.Errorf("failed to list pull request files: %w", err)
		}

		for _, file := range files {
			_, err := client.GetCommitFileContent(context.Background(), *file.Filename, *commit.SHA)
			if err != nil {
				return fmt.Errorf("failed to get commit file content: %w", err)
			}

			// 这里可以添加代码审查逻辑，例如使用AI进行审查
			// 示例：添加一个简单的评论
			comments = append(comments, &github.DraftReviewComment{
				Path:     github.String(*file.Filename),
				Body:     github.String("This is a sample comment."),
				Position: github.Int(1), // 评论的位置（行号）
			})
		}
	}

	// 提交代码审查
	err = client.SubmitCodeReview(context.Background(), *pendingCommits[0], comments)
	if err != nil {
		return fmt.Errorf("failed to submit code review: %w", err)
	}

	return nil
}
