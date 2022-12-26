package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/usual2970/userhub/domain"
	"github.com/usual2970/userhub/domain/constant"
)

const sqAccessTokenPre = "sq:access:token:%s"

type AccessTokenRepository struct {
	rc *redis.Client
}

func NewAccessTokenRepository(rc *redis.Client) domain.IAccessTokenRepository {
	return &AccessTokenRepository{
		rc: rc,
	}
}

func (atRepo *AccessTokenRepository) GetAccessToken(ctx context.Context, token string) (int, error) {
	key := atRepo.getKey(token)
	str, err := atRepo.rc.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}
func (atRepo *AccessTokenRepository) SetAccessToken(ctx context.Context, token string, id int) error {
	key := atRepo.getKey(token)
	return atRepo.rc.Set(ctx, key, id, constant.AuthExpireDuration).Err()
}

func (atRepo *AccessTokenRepository) DelAccessToken(ctx context.Context, token string) error {
	key := atRepo.getKey(token)
	return atRepo.rc.Del(ctx, key).Err()
}

func (atRepo *AccessTokenRepository) getKey(token string) string {
	return fmt.Sprintf(sqAccessTokenPre, token)
}
