package domain

import "context"

type IAccessTokenRepository interface {
	GetAccessToken(ctx context.Context, token string) (int, error)
	SetAccessToken(ctx context.Context, token string, id int) error
	DelAccessToken(ctx context.Context, token string) error
}
