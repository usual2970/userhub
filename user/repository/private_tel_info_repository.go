package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/usual2970/userhub/domain"
	"github.com/usual2970/userhub/domain/constant"
	"gorm.io/gorm"
)

const sqPrivateTelInfoPre = "sq:private:tel:info:%s"

type PrivateTelInfoRepository struct {
	db *gorm.DB
	rc *redis.Client
}

func NewPrivateTelInfoRepository(db *gorm.DB, rc *redis.Client) domain.IPrivateTelInfoRepository {
	return &PrivateTelInfoRepository{
		db: db,
		rc: rc,
	}
}

func (pTr *PrivateTelInfoRepository) GetOneByHash(ctx context.Context, hash string) (*domain.PrivateTelInfo, error) {
	key := pTr.getKey(hash)
	str, err := pTr.rc.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if str == constant.NotExistData {
		return nil, gorm.ErrRecordNotFound
	}

	rs := &domain.PrivateTelInfo{}
	if err == nil {
		if err := jsoniter.UnmarshalFromString(str, rs); err != nil {
			return nil, err
		}
		return rs, nil
	}

	if err := pTr.db.Where("tel_hash=?", hash).First(rs).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		pTr.rc.Set(ctx, key, constant.NotExistData, time.Hour*48)
		return nil, err
	}
	data, _ := jsoniter.MarshalToString(rs)
	pTr.rc.Set(ctx, key, data, time.Hour*48)
	return rs, nil
}

func (pTr *PrivateTelInfoRepository) DelByHashFromCache(ctx context.Context, hash string) error {
	key := pTr.getKey(hash)
	return pTr.rc.Del(ctx, key).Err()
}

func (pTr *PrivateTelInfoRepository) getKey(hash string) string {
	return fmt.Sprintf(sqPrivateTelInfoPre, hash)
}
