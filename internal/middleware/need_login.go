package middleware

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/usual2970/gopkg/conf"
	"github.com/usual2970/gopkg/container"
	pkgJwt "github.com/usual2970/gopkg/jwt"
	"github.com/usual2970/gopkg/log"
	"github.com/usual2970/userhub/domain"
)

func NeedLogin() echo.MiddlewareFunc {
	signingKey := []byte(conf.GetString("auth.key"))

	config := middleware.JWTConfig{
		TokenLookup: "header:Authorization:Bearer ",
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "HS256" {
					return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
				}
				return signingKey, nil
			}

			// claims are of type `jwt.MapClaims` when token is created with `jwt.Parse`
			token, err := jwt.Parse(auth, keyFunc)
			log.Info("token: %v", token.Raw, token.Valid)
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}

			var accessTokenRepo domain.IAccessTokenRepository
			container.Invoke(func(repo domain.IAccessTokenRepository) {
				accessTokenRepo = repo
			})

			if _, err := accessTokenRepo.GetAccessToken(c.Request().Context(), token.Raw); err != nil {
				return nil, errors.New("invalid token")
			}

			ctx := pkgJwt.WithContext(c.Request().Context(), token.Claims)

			c.SetRequest(c.Request().Clone(ctx))
			return token, nil
		},
	}

	return middleware.JWTWithConfig(config)
}
