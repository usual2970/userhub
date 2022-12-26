package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/usual2970/gopkg/conf"
	"github.com/usual2970/gopkg/log"
	_articleHttpDelivery "github.com/usual2970/userhub/article/delivery/http"
	_articleHttpDeliveryMiddleware "github.com/usual2970/userhub/article/delivery/http/middleware"
	_articleRepo "github.com/usual2970/userhub/article/repository/mysql"
	_articleUcase "github.com/usual2970/userhub/article/usecase"
	_authorRepo "github.com/usual2970/userhub/author/repository/mysql"

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

	dbHost := conf.GetString(`database.host`)
	dbPort := conf.GetString(`database.port`)
	dbUser := conf.GetString(`database.user`)
	dbPass := conf.GetString(`database.pass`)
	dbName := conf.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Error(err)
		return
	}
	err = dbConn.Ping()
	if err != nil {
		log.Error(err)
		return
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}()

	// database
	db, closeDb, err := registerDb()
	if err != nil {
		log.Error(err)
		return
	}
	defer func() {
		closeDb()
	}()

	redis, closeRedis, err := registerRedis()
	if err != nil {
		log.Error(err)
		return
	}
	defer func() {
		closeRedis()
	}()

	e := echo.New()
	middL := _articleHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)

	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	_articleHttpDelivery.NewArticleHandler(e, au)

	accessTokenRepo := userRepo.NewAccessTokenRepository(redis)

	accountRepo := userRepo.NewAccountRepository(db, redis)

	codeRepo := userRepo.NewCodeRepository(db, redis)

	telInfoRepo := userRepo.NewPrivateTelInfoRepository(db, redis)

	authUsecase := userUsecase.NewAuthUsecase(codeRepo, accountRepo, telInfoRepo, accessTokenRepo)

	userHttpDelivery.NewAuthHandler(e, authUsecase)

	for _, route := range e.Routes() {
		fmt.Println(route.Path)
	}
	log.Error(e.Start(viper.GetString("server.address"))) //nolint
}

func registerDb() (*gorm.DB, func() error, error) {

	dbHost := conf.GetString(`database.host`)
	dbPort := conf.GetString(`database.port`)
	dbUser := conf.GetString(`database.user`)
	dbPass := conf.GetString(`database.pass`)
	dbName := conf.GetString(`database.name`)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	close := func() error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return db, close, nil
}

func registerRedis() (*redis.Client, func() error, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	close := func() error {
		return rdb.Close()
	}

	return rdb, close, nil
}
