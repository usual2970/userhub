package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	goRedis "github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/usual2970/gopkg/gorm"
	"github.com/usual2970/gopkg/redis"
	"github.com/usual2970/userhub/domain"
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
}

func NewCodeRepository() domain.ICodeRepository {
	return &CodeRepository{}
}

func (cr *CodeRepository) Save(ctx context.Context, code *domain.Code) error {
	db, err := gorm.GetDB()
	if err != nil {
		return err
	}
	if err := db.Save(code).Error; err != nil {
		return err
	}
	cr.saveToCache(ctx, code)
	return nil
}

func (cr *CodeRepository) Update(ctx context.Context, code *domain.Code) error {
	db, err := gorm.GetDB()
	if err != nil {
		return err
	}
	if err := db.Save(code).Error; err != nil {
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
	rc, err := redis.GetRedis()
	if err != nil {
		return nil, err
	}
	key := getCodeKey(tel, nationCode, purpose)

	str, err := rc.Get(ctx, key).Result()

	if err != nil && !errors.Is(err, goRedis.Nil) {
		return nil, err
	}

	if err != nil && errors.Is(err, goRedis.Nil) {
		return nil, domain.ErrNotFound
	}
	rs := &domain.Code{}
	_ = jsoniter.UnmarshalFromString(str, rs)
	return rs, nil
}

func (cr *CodeRepository) saveToCache(ctx context.Context, code *domain.Code) error {
	rc, err := redis.GetRedis()
	if err != nil {
		return err
	}
	key := getCodeKey(code.Tel, code.NationCode, code.Purpose)
	rs, _ := jsoniter.MarshalToString(code)
	rc.Set(ctx, key, rs, time.Duration(code.ExpiredAt-time.Now().Unix())*time.Second)
	return nil
}

func getCodeKey(tel, nationCode string, purpose int) string {
	return fmt.Sprintf(sqCodePrefix, tel, nationCode, purpose)
}
