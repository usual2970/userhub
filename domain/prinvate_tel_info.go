package domain

import (
	"context"
	"time"
)

type PrivateTelInfo struct {
	ID        int    `gorm:"column:id;primaryKey"`
	TelHash   string `gorm:"column:tel_hash"`
	UserId    int    `gorm:"column:user_id"`
	CreatedAt int64  `gorm:"column:created_at"`
	DeletedAt int64  `gorm:"column:deleted_at"`
}

func (*PrivateTelInfo) TableName() string {
	return "user_private_tel_info"
}

func NewPrivateTelInfo(hash string) *PrivateTelInfo {
	now := time.Now().Unix()
	return &PrivateTelInfo{
		TelHash:   hash,
		CreatedAt: now,
	}
}

func (pti *PrivateTelInfo) SetUserID(Id int) {
	pti.UserId = Id
}

type IPrivateTelInfoRepository interface {
	GetOneByHash(ctx context.Context, hash string) (*PrivateTelInfo, error)
	DelByHashFromCache(ctx context.Context, hash string) error
}
