package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/usual2970/userhub/domain"
)

type ResponseError struct {
	Message string `json:"message"`
}

func Resp(ctx echo.Context, rs any) error {
	if rs == nil {
		rs = struct{}{}
	}
	return ctx.JSON(http.StatusOK, rs)
}

func Err(ctx echo.Context, err error) error {
	return ctx.JSON(getStatusCode(err), &ResponseError{Message: err.Error()})
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
