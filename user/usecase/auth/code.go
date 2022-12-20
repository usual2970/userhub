package auth

import (
	"context"
	"errors"
	"regexp"

	"github.com/usual2970/gopkg/log"
	tel2 "github.com/usual2970/gopkg/tel"
	"github.com/usual2970/userhub/domain"
	"github.com/usual2970/userhub/domain/constant"
	"github.com/usual2970/userhub/internal/openid"
	"gorm.io/gorm"
)

var ErrParamWrongTel = errors.New("param wrong tel")
var ErrParamWrongCode = errors.New("param wrong code")

var codeReg = regexp.MustCompile(`^\d{4}$`)

type CodeAuth struct {
	platformId  int
	codeRepo    domain.ICodeRepository
	accountRepo domain.IAccountRepository
	pTelRepo    domain.IPrivateTelInfoRepository
}

func NewCodeAuth(codeRepo domain.ICodeRepository, accountRepo domain.IAccountRepository, pTelRepo domain.IPrivateTelInfoRepository) *CodeAuth {
	return &CodeAuth{
		platformId:  constant.PlatformSmsCode,
		codeRepo:    codeRepo,
		accountRepo: accountRepo,
		pTelRepo:    pTelRepo,
	}
}

func (ca *CodeAuth) CheckParam(ctx context.Context, param map[string]string) error {
	tel, ok := param["tel"]
	if !ok || tel == "" {
		return ErrParamWrongTel
	}
	code, ok := param["code"]
	if !ok || !codeReg.MatchString(code) {
		return ErrParamWrongCode
	}

	nationCode := param["nationCode"]
	if nationCode == "" {
		nationCode = domain.DefaultNationCode
	}

	c, err := ca.codeRepo.GetByTelAndCode(ctx, tel, nationCode, code, domain.CodePurposeLogin)
	if err != nil {
		return ErrParamWrongCode
	}
	if c.IsUsed() {
		return ErrParamWrongCode
	}

	c.SetState(domain.CodeStateUsed)

	_ = ca.codeRepo.Update(ctx, c)
	return nil
}

func (ca *CodeAuth) Login(ctx context.Context, param map[string]string) (*domain.Account, error) {
	l := log.WithField("module", "sms code").WithField("param", param)
	// 先生成openid
	nationCode := param["nationCode"]
	tel := domain.DefaultNationCode + param["tel"]
	if nationCode != "" {
		tel = param["nationCode"] + param["tel"]
	}

	return telLogin(ctx, tel, ca.platformId, ca.accountRepo, ca.pTelRepo, l)
}

func telLogin(ctx context.Context, tel string, platformId int, accountRepo domain.IAccountRepository,
	pTelRepo domain.IPrivateTelInfoRepository, l *log.DLogger) (*domain.Account, error) {

	openId := openid.Openid(constant.PlatformSmsCode, tel)

	// 看openid是否已存在
	account, err := accountRepo.GetOneByOpenid(ctx, openId)
	l.Info(account, err)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// 存在则登录成功
	if err == nil {
		return account, nil
	}

	// 不存在则查看手机号是否存在
	telHash := tel2.Hash(tel)
	privateTel, err := pTelRepo.GetOneByHash(ctx, telHash)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	account = domain.NewAccount(openId, platformId)

	// 不存在则插入记录
	var privateInfo *domain.PrivateInfo

	var profile *domain.Profile

	if errors.Is(err, gorm.ErrRecordNotFound) {
		privateTel = domain.NewPrivateTelInfo(telHash)
		privateInfo = domain.NewPrivateInfo("", tel, telHash)
		profile = domain.InitProfile()
	}

	if err == nil {
		account.SetUserID(privateTel.UserId)
	}

	// 存在则生成openid插入记录
	if err := accountRepo.Save(ctx, account, privateTel, privateInfo, profile, func(account *domain.Account) {

		if account.UserId == 0 {
			account.SetUserID(account.ID)
		}

		if privateTel != nil {
			privateTel.SetUserID(account.ID)
		}
		if privateInfo != nil {
			privateInfo.SetUserID(account.ID)
		}
		if profile != nil {
			profile.SetUserID(account.ID)
		}
	}); err != nil {
		return nil, err
	}

	_ = accountRepo.DeleteFromCache(ctx, openId)
	return account, nil
}
