package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/usual2970/gopkg/conf"
	pkgJwt "github.com/usual2970/gopkg/jwt"
	"github.com/usual2970/gopkg/log"
	"github.com/usual2970/userhub/domain"
	"github.com/usual2970/userhub/domain/constant"
	"github.com/usual2970/userhub/user/usecase/auth"
)

const AuthExpireDuration = time.Hour * 24 * 7

var ErrAuthNotExist = errors.New("auth not found")
var ErrCodeHasSend = errors.New("code has send")

type IAuth interface {
	CheckParam(ctx context.Context, param map[string]string) error
	Login(ctx context.Context, param map[string]string) (*domain.Account, error)
}

type AuthUsecase struct {
	codeRepo    domain.ICodeRepository
	accountRepo domain.IAccountRepository
	pTelRepo    domain.IPrivateTelInfoRepository
	atRepo      domain.IAccessTokenRepository
}

func NewAuthUsecase(codeRepo domain.ICodeRepository, accountRepo domain.IAccountRepository,
	pTelRepo domain.IPrivateTelInfoRepository, atRepo domain.IAccessTokenRepository) domain.IAuthUsecase {
	return &AuthUsecase{
		codeRepo:    codeRepo,
		accountRepo: accountRepo,
		pTelRepo:    pTelRepo,
		atRepo:      atRepo,
	}
}

func (a *AuthUsecase) getAuth(platformId int) (IAuth, error) {
	switch platformId {
	case constant.PlatformSmsCode:
		return auth.NewCodeAuth(a.codeRepo, a.accountRepo, a.pTelRepo), nil
	default:
		return nil, ErrAuthNotExist
	}
}

// Login 登录
func (a *AuthUsecase) Login(ctx context.Context, param *domain.AuthLoginReq) (*domain.AuthLoginResp, error) {
	l := log.WithField("module", "login").WithField("param", param)

	l.Info("login start")
	// 获取登录对象
	auth, err := a.getAuth(param.PlatformId)
	if err != nil {
		l.Error("get auth failed:", err)
		return nil, err
	}
	// 验证参数
	if err := auth.CheckParam(ctx, param.Data); err != nil {
		l.Error("check param failed:", err)
		return nil, err
	}

	// 登录流程
	account, err := auth.Login(ctx, param.Data)
	if err != nil {
		l.Error("login failed:", err)
		return nil, err
	}
	// 生成jwt
	expiresAt := time.Now().Add(AuthExpireDuration)
	claims := &jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%d", account.UserId),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := conf.GetString("auth.key")
	if key == "" {
		l.Error("get auth key failed:")
		return nil, errors.New("get auth key failed")
	}
	accessToken, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}
	rs := &domain.AuthLoginResp{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt.Unix(),
	}
	err = a.atRepo.SetAccessToken(ctx, accessToken, account.UserId)
	if err != nil {
		l.Error("save accesstoken failed:", err)
		return nil, err
	}
	return rs, nil
}

// SmsCode 发送短信验证码
func (a *AuthUsecase) SmsCode(ctx context.Context, param *domain.AuthSmsCodeReq) error {

	l := log.WithField("module", "sms code").WithField("param", param)
	if param.NationCode == "" {
		param.NationCode = domain.DefaultNationCode
	}
	// 先查询有没有code
	if c, err := a.codeRepo.GetByTel(ctx, param.Tel, param.NationCode, domain.CodePurposeLogin); err != nil && !errors.Is(err, domain.ErrNotFound) {
		return errors.New("send sms code failed:" + err.Error())
	} else if err == nil && !c.IsUsed() {
		return ErrCodeHasSend
	}
	// 没有的话再发送
	code := domain.NewSmsCode(param.Tel, param.NationCode, domain.CodePurposeLogin)
	err := a.codeRepo.Save(ctx, code)
	if err != nil {
		l.Error("send sms code failed: ", err)
		return errors.New("send sms code failed:" + err.Error())
	}

	// 发送短信
	return nil
}

// Logout 登出
func (a *AuthUsecase) Logout(ctx context.Context) error {

	token, err := pkgJwt.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	return a.atRepo.DelAccessToken(ctx, token)
}
