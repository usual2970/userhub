package domain

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/rs/xid"
)

type Profile struct {
	UserId         int    `gorm:"column:user_id;primaryKey"`
	Uri            string `gorm:"column:uri"`
	Headimgurl     string `gorm:"column:headimgurl"`
	Signature      string `gorm:"column:signature"`
	Nickname       string `gorm:"column:nickname"`
	Sex            int    `gorm:"column:sex"`
	Region         string `gorm:"column:region"`
	Country        string `gorm:"column:country"`
	Province       string `gorm:"column:province"`
	City           string `gorm:"column:city"`
	Height         int    `gorm:"height"`
	Weight         int    `gorm:"weight"`
	Lang           string `gorm:"column:lang"`
	Preferences    string `gorm:"column:preferences"`
	Birthday       string `gorm:"column:birthday"`
	ForbiddenAt    int64  `gorm:"column:forbidden_at"`
	ForbiddenEndAt int64  `gorm:"column:forbidden_end_at"`
	CreatedAt      int64  `gorm:"column:created_at"`
	UpdatedAt      int64  `gorm:"column:updated_at"`
	DeletedAt      int64  `gorm:"column:deleted_at"`
}

func (p *Profile) TableName() string {
	return "user_profile"
}

func InitProfile() *Profile {
	now := time.Now().Unix()
	return &Profile{
		Uri:       xid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewProfile(snsInfo *WechatUserInfoResp) *Profile {
	now := time.Now().Unix()
	return &Profile{
		Uri:        xid.New().String(),
		Headimgurl: snsInfo.Headimgurl,
		Nickname:   snsInfo.Nickname,
		Sex:        snsInfo.Sex,
		Province:   snsInfo.Province,
		City:       snsInfo.City,
		Country:    snsInfo.Country,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (p *Profile) SetUserID(id int) {
	p.UserId = id
}

type IProfileRepository interface {
	UpdateByID(ctx context.Context, id int, data map[string]interface{}) error
	GetOneByID(ctx context.Context, id int) (*Profile, error)
}

type WechatAccessTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Unionid      string `json:"unionid"`
}

type WechatUserInfoResp struct {
	Openid     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Headimgurl string `json:"headimgurl"`
	Privilege  string `json:"privilege"`
	Unionid    string `json:"unionid"`
}

type WechatErrResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type IWechatRepository interface {
	GetAccessToken(code string) (*WechatAccessTokenResp, error)
	UserInfo(accessToken, openid string) (*WechatUserInfoResp, error)
}

type IProviceUsecase interface {
	// SetHead 设置头像
	SetHead(ctx context.Context, fh *multipart.FileHeader) (*UserSetHeadResp, error)
	// SetNickname 设置昵称
	SetNickname(ctx context.Context, nickname string) error
	// SetSex 设置性别
	SetSex(ctx context.Context, sex int) error
	// SetCity 设置城市
	SetCity(ctx context.Context, req *UserSetCityReq) error
	// Profile 获取用户基本信息
	Profile(ctx context.Context) (*UserProfileResp, error)
	// SetSignature 设置签名
	SetSignature(ctx context.Context, signature string) error
	// SetHeight 设置身高
	SetHeight(ctx context.Context, height int) error
	// SetWeight 设置体重
	SetWeight(ctx context.Context, weight int) error
	// SetBirthday 设置生日
	SetBirthday(ctx context.Context, birthday string) error
}
