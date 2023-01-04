package main

import (
	"fmt"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	"github.com/usual2970/gopkg/conf"
	"github.com/usual2970/gopkg/container"
	"github.com/usual2970/gopkg/gorm"
	"github.com/usual2970/gopkg/log"
	"github.com/usual2970/gopkg/redis"
	"github.com/usual2970/userhub/domain"

	userRepo "github.com/usual2970/userhub/user/repository"

	userHttpDelivery "github.com/usual2970/userhub/user/delivery/http"
	userUsecase "github.com/usual2970/userhub/user/usecase"
)

func init() {
	// 日志
	log.Setup()
	// 配置
	absPath, err := filepath.Abs("..")
	if err != nil {
		panic(err)
	}
	if err := conf.Setup(conf.WithPath(absPath)); err != nil {
		panic(err)
	}

}

func main() {

	// database
	if err := registerCommon(); err != nil {
		log.Error(err)
		return
	}

	if err := registerUser(); err != nil {
		log.Error(err)
		return
	}

	var e *echo.Echo

	if err := container.Invoke(func(echo *echo.Echo) {
		e = echo
	}); err != nil {
		log.Error(err)
		return
	}

	db, err := gorm.GetDB()
	if err != nil {
		log.Error(err)
		return
	}

	red, err := redis.GetRedis()
	if err != nil {
		log.Error(err)
		return
	}

	defer func() {
		sqlDb, err := db.DB()
		if err == nil {
			sqlDb.Close()
		}

		red.Close()
	}()

	for _, route := range e.Routes() {
		fmt.Println(route.Path)
	}
	log.Error(e.Start(viper.GetString("server.address"))) //nolint
}

func registerUser() error {
	if err := container.Provide(userRepo.NewAccessTokenRepository); err != nil {
		return err
	}

	if err := container.Provide(userRepo.NewAccountRepository); err != nil {
		return err
	}

	if err := container.Provide(userRepo.NewCodeRepository); err != nil {
		return err
	}

	if err := container.Provide(userRepo.NewPrivateTelInfoRepository); err != nil {
		return err
	}

	if err := container.Provide(userUsecase.NewAuthUsecase); err != nil {
		return err
	}

	if err := container.Invoke(func(e *echo.Echo, authUc domain.IAuthUsecase) {
		userHttpDelivery.NewAuthHandler(e, authUc)
	}); err != nil {
		return err
	}

	return nil
}

func registerCommon() error {

	// provde echo
	if err := container.Provide(func() *echo.Echo {
		e := echo.New()
		e.Use(middleware.CORS())

		timeout := conf.GetInt("context.timeout")

		e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Skipper: middleware.DefaultSkipper,
			Timeout: time.Duration(timeout) * time.Second,
		}))

		e.Use(middleware.Logger())
		return e
	}); err != nil {
		return err
	}
	return nil
}
