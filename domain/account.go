package domain

import (
	"context"
	"time"
)

type Account struct {
	ID         int    `gorm:"column:id;primaryKey"`
	Openid     string `gorm:"column:openid"`
	PlatformId int    `gorm:"column:platform_id"`
	UserId     int    `gorm:"column:user_id"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdatedAt  int64  `gorm:"column:updated_at"`
	DeletedAt  int64  `gorm:"column:deleted_at"`
}

func (*Account) TableName() string {
	return "user_account"
}

func NewAccount(openid string, platformId int) *Account {
	now := time.Now().Unix()
	return &Account{
		Openid:     openid,
		PlatformId: platformId,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (a *Account) SetUserID(id int) {
	a.UserId = id
	a.UpdatedAt = time.Now().Unix()
}

type IAccountRepository interface {
	GetOneByOpenid(ctx context.Context, openid string) (*Account, error)

	GetByUserID(ctx context.Context, userId int) ([]Account, error)
	DelByUserIdFromCache(ctx context.Context, userId int) error

	Save(ctx context.Context, account *Account, privateTet *PrivateTelInfo, privateInfo *PrivateInfo, profile *Profile, update func(account *Account)) error

	SaveWechat(ctx context.Context, account *Account, profile *Profile, unionid *Unionid, update func(account *Account)) error

	DeleteFromCache(ctx context.Context, openid string) error
}

type IAccountUsecase interface {
	// Info 账户信息手机号、微信
	Info(ctx context.Context) (*AccountInfoResp, error)
	// BindTelCode 绑定手机号验证码
	BindTelCode(ctx context.Context, param *AccountBindTelCodeReq) error
	// BindTel 绑定手机号
	BindTel(ctx context.Context, param *AccountBindTelReq) error
	// BindWechat 绑定微信
	BindWechat(ctx context.Context, param *AccountBindWechatReq) error
}

type IAuthUsecase interface {
	// Login 登录
	Login(ctx context.Context, param *AuthLoginReq) (*AuthLoginResp, error)
	// SmsCode 发送短信验证码
	SmsCode(ctx context.Context, param *AuthSmsCodeReq) error
	// Logout 登出
	Logout(ctx context.Context) error
}
