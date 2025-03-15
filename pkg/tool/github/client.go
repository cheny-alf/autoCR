package github

import "github.com/google/go-github/v69/github"

const (
	X_GitHub_Api_Version = "2022-11-28"
	Accept               = "application/vnd.github+json"
)

func NewGithubClient(token string) *github.Client {
	return github.NewClient(nil).WithAuthToken(token)
}

// 根据pr获取相关地址
func GetCommitDetail(url string) {

}

// 创建评论
func CreateCommitComment() {

}
