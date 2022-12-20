package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

	"github.com/usual2970/gopkg/conf"
	"github.com/usual2970/gopkg/log"
	_articleHttpDelivery "github.com/usual2970/userhub/article/delivery/http"
	_articleHttpDeliveryMiddleware "github.com/usual2970/userhub/article/delivery/http/middleware"
	_articleRepo "github.com/usual2970/userhub/article/repository/mysql"
	_articleUcase "github.com/usual2970/userhub/article/usecase"
	_authorRepo "github.com/usual2970/userhub/author/repository/mysql"
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
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
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

	e := echo.New()
	middL := _articleHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	_articleHttpDelivery.NewArticleHandler(e, au)

	log.Error(e.Start(viper.GetString("server.address"))) //nolint
}
