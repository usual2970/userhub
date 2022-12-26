package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/usual2970/userhub/domain"
	"gorm.io/gorm"
)

const (
	sqCodePrefix = "sq:code:prefix:%s:%s:%d"
)

const getCodeScript = `
local e = redis.call("exists", KEYS[1])
if tonumber(e) == 0 then
    return nil
end
return redis.call("hmget", KEYS[1], "id", "code", "state")
`

type CodeRepository struct {
	db *gorm.DB
	rc *redis.Client
}

func NewCodeRepository(db *gorm.DB, rc *redis.Client) domain.ICodeRepository {
	return &CodeRepository{
		db: db,
		rc: rc,
	}
}

func (cr *CodeRepository) Save(ctx context.Context, code *domain.Code) error {
	if err := cr.db.Save(code).Error; err != nil {
		return err
	}
	cr.saveToCache(ctx, code)
	return nil
}

func (cr *CodeRepository) Update(ctx context.Context, code *domain.Code) error {
	if err := cr.db.Save(code).Error; err != nil {
		return err
	}

	cr.saveToCache(ctx, code)
	return nil
}

func (cr *CodeRepository) GetByTelAndCode(ctx context.Context, tel, nationCode, code string, purpose int) (*domain.Code, error) {
	rs, err := cr.GetByTel(ctx, tel, nationCode, purpose)
	if err != nil {
		return nil, err
	}

	if rs.Code != code {
		return nil, domain.ErrNotFound
	}
	return rs, nil
}

func (cr *CodeRepository) GetByTel(ctx context.Context, tel, nationCode string, purpose int) (*domain.Code, error) {
	key := getCodeKey(tel, nationCode, purpose)

	str, err := cr.rc.Get(ctx, key).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if err != nil && errors.Is(err, redis.Nil) {
		return nil, domain.ErrNotFound
	}
	rs := &domain.Code{}
	_ = jsoniter.UnmarshalFromString(str, rs)
	return rs, nil
}

func (cr *CodeRepository) saveToCache(ctx context.Context, code *domain.Code) {
	key := getCodeKey(code.Tel, code.NationCode, code.Purpose)
	rs, _ := jsoniter.MarshalToString(code)
	cr.rc.Set(ctx, key, rs, time.Duration(code.ExpiredAt-time.Now().Unix())*time.Second)
}

func getCodeKey(tel, nationCode string, purpose int) string {
	return fmt.Sprintf(sqCodePrefix, tel, nationCode, purpose)
}
