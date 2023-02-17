package auth

import (
	"context"
	"errors"

	"github.com/usual2970/gopkg/log"

	"github.com/usual2970/userhub/domain"
	"github.com/usual2970/userhub/domain/constant"
	"github.com/usual2970/userhub/internal/openid"
	"gorm.io/gorm"
)

type GithubAuth struct {
	githubRepo  domain.IGithubRepository
	accountRepo domain.IAccountRepository
}

func NewGithubAuth(githubRepo domain.IGithubRepository, accountRepo domain.IAccountRepository) *GithubAuth {
	return &GithubAuth{
		githubRepo:  githubRepo,
		accountRepo: accountRepo,
	}
}

func (ga *GithubAuth) CheckParam(ctx context.Context, param map[string]string) error {
	if _, ok := param["code"]; !ok {
		return ErrParamWrongCode
	}
	return nil
}
func (ga *GithubAuth) Login(ctx context.Context, param map[string]string) (*domain.Account, error) {
	l := log.WithField("module", "sms_code").WithField("param", param)
	// 获取accessToken
	tokenRs, err := ga.githubRepo.GetAccessToken(param["code"])
	if err != nil {
		l.Info("get access token failed:", tokenRs)
		return nil, err
	}

	// 获取snsInfo
	ghInfo, err := ga.githubRepo.UserInfo(tokenRs.AccessToken, "")
	if err != nil {
		return nil, err
	}
	openId := openid.Openid(constant.PlatformGithub, ghInfo.Openid)
	// account是否存在
	account, err := ga.accountRepo.GetOneByOpenid(ctx, openId)
	l.Info(account, err)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// 存在则登录成功
	if err == nil {
		return account, nil
	}

	account = domain.NewAccount(openId, constant.PlatformGithub)

	profile := domain.NewProfileWithGithub(ghInfo)

	if err := ga.accountRepo.SaveSns(ctx, account, profile, nil, func(account *domain.Account) {
		if account.UserId == 0 {
			account.SetUserID(account.ID)
		}

		if profile != nil {
			profile.SetUserID(account.ID)
		}
	}); err != nil {
		return nil, err
	}

	return account, nil
}
