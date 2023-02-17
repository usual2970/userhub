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
	"github.com/usual2970/gopkg/sql"
	"github.com/usual2970/userhub/domain"
	"github.com/usual2970/userhub/domain/constant"
	goGorm "gorm.io/gorm"
)

const sqAccountPrefix = "sq:account:%s"
const sqAccountsPrefix = "sq:accounts:%d"

type AccountRepository struct {
}

func NewAccountRepository() domain.IAccountRepository {
	return &AccountRepository{}
}

func (ar *AccountRepository) GetOneByOpenid(ctx context.Context, openid string) (*domain.Account, error) {
	rc, err := redis.GetRedis()
	if err != nil {
		return nil, err
	}
	key := ar.getAccountKey(openid)
	str, err := rc.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, goRedis.Nil) {
		return nil, err
	}

	if str == constant.NotExistData {
		return nil, goGorm.ErrRecordNotFound
	}

	rs := &domain.Account{}
	if err == nil {
		if err := jsoniter.UnmarshalFromString(str, rs); err != nil {
			return nil, err
		}
		return rs, nil
	}

	db, err := gorm.GetDB()

	if err := db.Where("openid=? and deleted_at=0", openid).First(rs).Error; err != nil {
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

func (ar *AccountRepository) Save(
	_ context.Context,
	account *domain.Account,
	privateTel *domain.PrivateTelInfo,
	privateInfo *domain.PrivateInfo,
	profile *domain.Profile,
	update func(account *domain.Account),
) error {
	db, dberr := gorm.GetDB()
	if dberr != nil {
		return dberr
	}
	tx := db.Begin()
	var err error
	defer func() {
		err = sql.FinishTransaction(err, tx)
	}()
	err = tx.Save(account).Error
	if err != nil {
		return err
	}

	updateFlag := false
	if account.UserId == 0 {
		updateFlag = true
	}

	update(account)

	if updateFlag {
		err = tx.Save(account).Error
		if err != nil {
			return err
		}
	}

	if privateTel != nil {
		err = tx.Save(privateTel).Error
		if err != nil {
			return err
		}
	}

	if privateInfo != nil {
		err = tx.Save(privateInfo).Error
		if err != nil {
			return err
		}
	}

	if profile != nil {
		err = tx.Save(profile).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (ar *AccountRepository) SaveSns(ctx context.Context, account *domain.Account, profile *domain.Profile, unionid *domain.Unionid, update func(account *domain.Account)) error {
	db, dberr := gorm.GetDB()
	if dberr != nil {
		return dberr
	}
	tx := db.Begin()
	var err error
	defer func() {
		err = sql.FinishTransaction(err, tx)
	}()
	err = tx.Save(account).Error
	if err != nil {
		return err
	}

	updateFlag := false
	if account.UserId == 0 {
		updateFlag = true
	}

	update(account)

	if updateFlag {
		err = tx.Save(account).Error
		if err != nil {
			return err
		}
	}

	if profile != nil {
		err = tx.Save(profile).Error
		if err != nil {
			return err
		}
	}

	if unionid != nil {
		err = tx.Save(unionid).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (ar *AccountRepository) DeleteFromCache(ctx context.Context, openid string) error {
	rc, err := redis.GetRedis()
	if err != nil {
		return err
	}
	key := ar.getAccountKey(openid)
	return rc.Del(ctx, key).Err()
}

func (ar *AccountRepository) GetByUserID(ctx context.Context, userId int) ([]domain.Account, error) {
	rc, err := redis.GetRedis()
	if err != nil {
		return nil, err
	}
	key := ar.getAccountsKey(userId)
	str, err := rc.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, goRedis.Nil) {
		return nil, err
	}

	rs := make([]domain.Account, 0)
	if err == nil {
		if err := jsoniter.UnmarshalFromString(str, rs); err != nil {
			return nil, err
		}
		return rs, nil
	}

	db, dberr := gorm.GetDB()
	if dberr != nil {
		return nil, dberr
	}

	if err := db.Where("user_id=? and deleted_at=0", userId).Find(&rs).Error; err != nil {
		return nil, err
	}
	data, _ := jsoniter.MarshalToString(rs)
	rc.Set(ctx, key, data, time.Hour*48)

	return rs, nil
}
func (ar *AccountRepository) DelByUserIdFromCache(ctx context.Context, userId int) error {
	rc, err := redis.GetRedis()
	if err != nil {
		return err
	}
	key := ar.getAccountsKey(userId)
	return rc.Del(ctx, key).Err()
}

func (ar *AccountRepository) getAccountKey(openid string) string {
	return fmt.Sprintf(sqAccountPrefix, openid)
}

func (ar *AccountRepository) getAccountsKey(userID int) string {
	return fmt.Sprintf(sqAccountsPrefix, userID)
}
