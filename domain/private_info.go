package domain

import (
	"context"
	"time"
)

type PrivateInfo struct {
	UserId       int    `gorm:"column:user_id;primaryKey"`
	Name         string `gorm:"column:name"`
	Tel          string `gorm:"column:tel"`
	TelHash      string `gorm:"column:tel_hash"`
	IdCode       string `gorm:"column:id_code"`
	RealInfo     string `gorm:"column:real_info"`
	IdCodeHash   string `gorm:"column:id_code_hash"`
	IdType       string `gorm:"column:id_type"`
	IdHold       string `gorm:"column:id_hold"`
	IdFront      string `gorm:"column:id_front"`
	IdBack       string `gorm:"column:id_back"`
	IdVerifiedAt int64  `gorm:"column:id_verified_at"`
	IdExpiredAt  int64  `gorm:"column:id_expired_at"`
	CreatedAt    int64  `gorm:"column:created_at`
}

func (*PrivateInfo) TableName() string {
	return "user_private_info"
}

func NewPrivateInfo(name, tel, telHash string) *PrivateInfo {
	return &PrivateInfo{
		Name:      name,
		Tel:       tel,
		TelHash:   telHash,
		CreatedAt: time.Now().Unix(),
	}
}

func (pt *PrivateInfo) SetUserID(Id int) {
	pt.UserId = Id
}

func (pt *PrivateInfo) HasTel() bool {
	return pt.Tel != ""
}

type IPrivateInfoRepository interface {
	GetOneByUserID(ctx context.Context, userID int) (*PrivateInfo, error)
	DetByUserIDFromCache(ctx context.Context, userID int) error
}
