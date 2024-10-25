package github

import (
	"context"
	"github.com/GitEval/GitEval-backend/internal/conf"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type gitHubAPI struct {
	client *github.Client
}

type GitHubAPI interface {
	GetLoginUrl() string
	SetClient(code string) error
	GetUserInfo(ctx context.Context) (*github.User, error)
}

// 创建一个全局的配置,由于每个账户都需要读取所以这个地方需要做这个事情,但是配置又需要放在configs中统一使用proto管理
var GitHubCfg struct {
	clientID     string
	clientSecret string
}

func InitGitHubCfg(cfg *conf.GitHub) {
	GitHubCfg.clientID = cfg.ClientID
	GitHubCfg.clientSecret = cfg.ClientSecret
}

func NewGitHubAPI() GitHubAPI {
	return &gitHubAPI{}
}

func (g *gitHubAPI) GetLoginUrl() string {
	redirectURL := "https://github.com/login/oauth/authorize?client_id=" + GitHubCfg.clientID + "&scope=user"
	return redirectURL
}

func (g *gitHubAPI) SetClient(code string) error {
	// 获取 access token
	token, err := g.getAccessToken(code)
	if err != nil {
		return err
	}

	// 使用 access token 创建 GitHub 客户端
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	g.client = client
	return nil
}

func (g *gitHubAPI) GetUserInfo(ctx context.Context) (*github.User, error) {
	// 获取用户信息
	user, _, err := g.client.Users.Get(ctx, "")
	if err != nil {
		return &github.User{}, err
	}
	return user, nil
}

func (g *gitHubAPI) getAccessToken(code string) (string, error) {
	// 创建 OAuth2 端点
	oauth2Endpoint := oauth2.Endpoint{
		TokenURL: "https://github.com/login/oauth/access_token",
	}

	// 创建 OAuth2 客户端
	ctx := context.Background()
	cf := oauth2.Config{
		ClientID:     GitHubCfg.clientID,
		ClientSecret: GitHubCfg.clientSecret,
		Endpoint:     oauth2Endpoint,
	}

	// 获取访问令牌
	token, err := cf.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
