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
	"github.com/usual2970/userhub/domain/constant"
	goGorm "gorm.io/gorm"
)

const sqPrivateTelInfoPre = "sq:private:tel:info:%s"

type PrivateTelInfoRepository struct {
}

func NewPrivateTelInfoRepository() domain.IPrivateTelInfoRepository {
	return &PrivateTelInfoRepository{}
}

func (pTr *PrivateTelInfoRepository) GetOneByHash(ctx context.Context, hash string) (*domain.PrivateTelInfo, error) {
	rc, err := redis.GetRedis()
	if err != nil {
		return nil, err
	}
	key := pTr.getKey(hash)
	str, err := rc.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, goRedis.Nil) {
		return nil, err
	}

	if str == constant.NotExistData {
		return nil, goGorm.ErrRecordNotFound
	}

	rs := &domain.PrivateTelInfo{}
	if err == nil {
		if err := jsoniter.UnmarshalFromString(str, rs); err != nil {
			return nil, err
		}
		return rs, nil
	}

	db, err := gorm.GetDB()
	if err != nil {
		return nil, err
	}

	if err := db.Where("tel_hash=?", hash).First(rs).Error; err != nil {
		if !errors.Is(err, goGorm.ErrRecordNotFound) {
			return nil, err
		}
		rc.Set(ctx, key, constant.NotExistData, time.Hour*48)
		return nil, err
	}
	data, _ := jsoniter.MarshalToString(rs)
	rc.Set(ctx, key, data, time.Hour*48)
	return rs, nil
}

func (pTr *PrivateTelInfoRepository) DelByHashFromCache(ctx context.Context, hash string) error {
	rc, err := redis.GetRedis()
	if err != nil {
		return err
	}
	key := pTr.getKey(hash)
	return rc.Del(ctx, key).Err()
}

func (pTr *PrivateTelInfoRepository) getKey(hash string) string {
	return fmt.Sprintf(sqPrivateTelInfoPre, hash)
}
