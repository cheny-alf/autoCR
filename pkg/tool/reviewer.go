package tool

import (
	"context"

	"github.com/google/go-github/github"
)

type Reviewer interface {
	GetCommitFileContent(ctx context.Context, filePath, sha string) (string, error)
	GetPendingReviewCommits() ([]*github.RepositoryCommit, error)
	GetPullRequestReviews(ctx context.Context) ([]*github.PullRequestReview, error)
	SubmitCodeReview(ctx context.Context, commit github.RepositoryCommit, comments []*github.DraftReviewComment) error
}

// 执行入口
func CodeReview(url string) {
	// 拿到客户端
	// 解析url
	// 获取pr内容
	// review
	// 提交

}
