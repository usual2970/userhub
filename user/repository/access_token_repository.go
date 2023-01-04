package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/usual2970/gopkg/redis"
	"github.com/usual2970/userhub/domain"
	"github.com/usual2970/userhub/domain/constant"
)

const sqAccessTokenPre = "sq:access:token:%s"

type AccessTokenRepository struct {
}

func NewAccessTokenRepository() domain.IAccessTokenRepository {
	return &AccessTokenRepository{}
}

func (atRepo *AccessTokenRepository) GetAccessToken(ctx context.Context, token string) (int, error) {
	rc, err := redis.GetRedis()
	if err != nil {
		return 0, err
	}
	key := atRepo.getKey(token)
	str, err := rc.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}
func (atRepo *AccessTokenRepository) SetAccessToken(ctx context.Context, token string, id int) error {
	rc, err := redis.GetRedis()
	if err != nil {
		return err
	}
	key := atRepo.getKey(token)
	return rc.Set(ctx, key, id, constant.AuthExpireDuration).Err()
}

func (atRepo *AccessTokenRepository) DelAccessToken(ctx context.Context, token string) error {
	rc, err := redis.GetRedis()
	if err != nil {
		return err
	}
	key := atRepo.getKey(token)
	return rc.Del(ctx, key).Err()
}

func (atRepo *AccessTokenRepository) getKey(token string) string {
	return fmt.Sprintf(sqAccessTokenPre, token)
}
