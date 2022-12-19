package domain

import (
	"context"
	"time"

	"jiqi/internal/common/code"
)

const expireDuration = time.Minute * 5
const (
	_ = iota
	CodePurposeLogin
	CodePurposeBindTel
)

const (
	_ = iota
	CodeStateSent
	CodeStateUsed
)

const DefaultNationCode = "0086"

type Code struct {
	ID         int    `gorm:"column:id;primaryKey"`
	Tel        string `gorm:"column:tel"`
	NationCode string `gorm:"column:nation_code"`
	Code       string `gorm:"column:code"`
	State      int    `gorm:"column:state"`
	Purpose    int    `gorm:"column:purpose"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdateAt   int64  `gorm:"column:updated_at"`
	ExpiredAt  int64  `gorm:"column:expired_at"`
}

func (*Code) TableName() string {
	return "user_code"
}

func NewSmsCode(tel, nationCode string, purpose int) *Code {
	now := time.Now()
	smsCode := code.Sms()
	return &Code{
		Tel:        tel,
		NationCode: nationCode,
		Purpose:    purpose,
		Code:       smsCode,
		State:      CodeStateSent,
		CreatedAt:  now.Unix(),
		UpdateAt:   now.Unix(),
		ExpiredAt:  now.Add(expireDuration).Unix(),
	}
}

func (c *Code) SetState(state int) {
	c.State = state
	c.UpdateAt = time.Now().Unix()
}

func (c *Code) IsUsed() bool {
	return c.State == CodeStateUsed
}

type ICodeRepository interface {
	Save(ctx context.Context, code *Code) error
	GetByTelAndCode(ctx context.Context, tel, nationCode, code string, purpose int) (*Code, error)
	GetByTel(ctx context.Context, tel, nationCode string, purpose int) (*Code, error)
	Update(ctx context.Context, code *Code) error
}
