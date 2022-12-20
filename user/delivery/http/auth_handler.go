package http

import (
	"github.com/labstack/echo"
	"github.com/usual2970/userhub/domain"
	"github.com/usual2970/userhub/internal/http"
)

type AuthHandler struct {
	AUsercase domain.IAuthUsecase
}

func NewAuthHandler(e *echo.Echo, aUsercase domain.IAuthUsecase) {
	handler := &AuthHandler{
		AUsercase: aUsercase,
	}
	g := e.Group("/user/v1/")
	g.POST("/auth/login", handler.Login)
	g.POST("/auth/sms-code", handler.SmsCode)
	g.POST("/auth/logout", handler.Logout)
}

func (a *AuthHandler) Login(ctx echo.Context) error {
	param := &domain.AuthLoginReq{}
	if err := ctx.Bind(param); err != nil {
		return http.Err(ctx, err)
	}
	resp, err := a.AUsercase.Login(ctx.Request().Context(), param)
	if err != nil {
		return http.Err(ctx, err)
	}
	return http.Resp(ctx, resp)
}

func (a *AuthHandler) SmsCode(ctx echo.Context) error {
	param := &domain.AuthSmsCodeReq{}
	if err := ctx.Bind(param); err != nil {
		return http.Err(ctx, err)
	}
	err := a.AUsercase.SmsCode(ctx.Request().Context(), param)
	if err != nil {
		return http.Err(ctx, err)
	}
	return http.Resp(ctx, nil)
}

func (a *AuthHandler) Logout(ctx echo.Context) error {
	err := a.AUsercase.Logout(ctx.Request().Context())
	if err != nil {
		return http.Err(ctx, err)
	}
	return http.Resp(ctx, nil)
}
