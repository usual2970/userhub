package domain

import (
	"context"
	"time"
)

const (
	_ = iota
	UnionidDtWechat
)

type Unionid struct {
	ID        int    `gorm:"column:id;primaryKey"`
	UserId    int    `gorm:"column:user_id"`
	Unionid   string `gorm:"column:unionid"`
	DataType  int    `gorm:"column:data_type"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func NewUnionid(unionid string, dt int) *Unionid {
	now := time.Now().Unix()
	return &Unionid{
		Unionid:   unionid,
		DataType:  dt,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (*Unionid) TableName() string {
	return "user_unionid"
}

func (u *Unionid) SetUserID(id int) {
	u.UserId = id
}

type IUnionidRepository interface {
	GetOneByUnionIdAndType(ctx context.Context, unionid string, dt int) (*Unionid, error)
	DelByUnionIdAndTypeFromCache(ctx context.Context, unionid string, dt int) error
}
