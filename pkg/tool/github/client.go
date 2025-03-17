package github

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v69/github"
)

const (
	X_GitHub_Api_Version = "2022-11-28"
	Accept               = "application/vnd.github+json"
)

func getToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

type GithubClient struct {
	Client      *github.Client
	Owner, Repo string
}

func NewGithubClient() *GithubClient {
	return &GithubClient{
		Client: github.NewClient(nil).WithAuthToken(getToken()),
	}
}

func (g *GithubClient) SetUserInfo(owner, repo string) *GithubClient {
	g.Owner = owner
	g.Repo = repo
	return g
}

// 根据pr获取相关地址
// getPRDetail retrieves detailed information about a Pull Request given its URL.
func getPRDetail(url string) (*github.PullRequest, error) {
	// Parse the URL to extract the owner, repo, and PR number.
	// Example: https://github.com/owner/repo/pull/123
	parts := strings.Split(url, "/")
	if len(parts) < 7 {
		return nil, fmt.Errorf("invalid URL format")
	}
	owner := parts[3]
	repo := parts[4]
	prNumberStr := parts[6]
	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PR number: %w", err)
	}

	// Fetch the Pull Request details.
	pr, _, err := NewGithubClient().SetUserInfo(owner, repo).Client.PullRequests.Get(context.Background(), owner, repo, prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR details: %w", err)
	}
	return pr, nil
}

// checkRepo verifies if the repository is accessible and valid.
func checkRepo(ctx context.Context, owner, repo string) error {
	client := NewGithubClient().SetUserInfo(owner, repo).Client
	_, resp, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to check repository: %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("repository not found or inaccessible")
	}
	return nil
}

// 获取指定commit中 文件内容
// GetCommitFileContent retrieves the content of a file at a specific commit and returns it with line numbers.
func (g *GithubClient) GetCommitFileContent(ctx context.Context, filePath, sha string) (string, error) {
	fileContent, _, _, err := g.Client.Repositories.GetContents(ctx, g.Owner, g.Repo, filePath, &github.RepositoryContentGetOptions{Ref: sha})
	if err != nil {
		return "", fmt.Errorf("failed to get file content: %w", err)
	}
	content, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("failed to decode file content: %w", err)
	}

	var contentWithLine []string
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		contentWithLine = append(contentWithLine, fmt.Sprintf("%d %s", i+1, line))
	}

	return strings.Join(contentWithLine, "\n"), nil
}

// GetPendingReviewCommits retrieves the commits that need review for a specific Pull Request.
func (g *GithubClient) GetPendingReviewCommits(ctx context.Context, prId int) ([]*github.RepositoryCommit, error) {
	commits, _, err := g.Client.PullRequests.ListCommits(ctx, g.Owner, g.Repo, prId, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list commits: %w", err)
	}

	var pendingReviewCommits []*github.RepositoryCommit
	for _, commit := range commits {
		if strings.HasPrefix(*commit.Commit.Message, "Merge pull request") {
			continue
		}

		commitInfo, _, err := g.Client.Repositories.GetCommit(ctx, g.Owner, g.Repo, *commit.SHA, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get commit info: %w", err)
		}
		pendingReviewCommits = append(pendingReviewCommits, commitInfo)
	}

	return pendingReviewCommits, nil
}

// GetPullRequestFiles retrieves the files associated with a specific Pull Request and prints their content.
func (g *GithubClient) GetPullRequestFiles(ctx context.Context, pullNum int) error {
	fileList, _, err := g.Client.PullRequests.ListFiles(ctx, g.Owner, g.Repo, pullNum, nil)
	if err != nil {
		return fmt.Errorf("failed to list PR files: %w", err)
	}

	for _, v := range fileList {
		fmt.Println(*v.SHA)
		fileContent, _, _, err := g.Client.Repositories.GetContents(ctx, g.Owner, g.Repo, *v.Filename, nil)
		if err != nil {
			return fmt.Errorf("failed to get file content: %w", err)
		}
		content, err := fileContent.GetContent()
		if err != nil {
			return fmt.Errorf("failed to decode file content: %w", err)
		}
		fmt.Print(content)
	}

	return nil
}

// 创建评论
// SubmitCodeReview submits a code review for a specific commit in a Pull Request.
// SubmitCodeReview submits a code review for a specific commit in a Pull Request.
func (g *GithubClient) SubmitCodeReview(ctx context.Context, commit github.RepositoryCommit, comments []*github.DraftReviewComment) error {
	prID, err := strconv.Atoi(*commit.Commit.Message)
	if err != nil {
		return fmt.Errorf("failed to parse PR ID from commit message: %w", err)
	}

	review, _, err := g.Client.PullRequests.CreateReview(ctx, g.Owner, g.Repo, prID, &github.PullRequestReviewRequest{
		Body:     github.String("Sample review"),
		Event:    github.String("COMMENT"),
		Comments: comments,
	})
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	log.Printf("Created review with ID: %d\n", *review.ID)
	return nil
}
