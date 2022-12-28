package main

import (
	"fmt"
	"path/filepath"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/usual2970/gopkg/conf"
	"github.com/usual2970/gopkg/container"
	"github.com/usual2970/gopkg/log"
	_articleHttpDeliveryMiddleware "github.com/usual2970/userhub/article/delivery/http/middleware"
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
	var db *gorm.DB
	var red *redis.Client

	if err := container.Invoke(func(echo *echo.Echo, gorm *gorm.DB, redis *redis.Client) {
		e = echo
		db = gorm
		red = redis
	}); err != nil {
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

func registerDb() (*gorm.DB, error) {

	dbHost := conf.GetString(`database.host`)
	dbPort := conf.GetString(`database.port`)
	dbUser := conf.GetString(`database.user`)
	dbPass := conf.GetString(`database.pass`)
	dbName := conf.GetString(`database.name`)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func registerRedis() (*redis.Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb, nil
}

func registerCommon() error {

	// provide db
	if err := container.Provide(registerDb); err != nil {
		return err
	}

	// provide redis
	if err := container.Provide(registerRedis); err != nil {
		return err
	}

	// provde echo
	if err := container.Provide(func() *echo.Echo {
		e := echo.New()
		middL := _articleHttpDeliveryMiddleware.InitMiddleware()
		e.Use(middL.CORS)

		e.Use(middleware.Logger())
		return e
	}); err != nil {
		return err
	}
	return nil
}
