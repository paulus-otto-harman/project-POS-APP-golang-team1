package infra

import (
	"project/config"
	"project/database"
	"project/handler"
	"project/infra/jwt"
	"project/log"
	"project/middleware"
	"project/repository"
	"project/service"

	"go.uber.org/zap"
)

type ServiceContext struct {
	Cacher     database.Cacher
	Cfg        config.Config
	Ctl        handler.Handler
	Log        *zap.Logger
	Middleware middleware.Middleware
	JWT        jwt.JWT
}

func NewServiceContext() (*ServiceContext, error) {

	handlerError := func(err error) (*ServiceContext, error) {
		return nil, err
	}

	// instance config
	appConfig, err := config.LoadConfig()
	if err != nil {
		return handlerError(err)
	}

	// instance logger
	logger, err := log.InitZapLogger(appConfig)
	if err != nil {
		return handlerError(err)
	}

	// instance database
	db, err := database.ConnectDB(appConfig)
	if err != nil {
		return handlerError(err)
	}

	jwtLib := jwt.NewJWT(appConfig.PrivateKey, appConfig.PublicKey, logger)

	rdb := database.NewCacher(appConfig, 60*60)

	// instance repository
	repo := repository.NewRepository(db, rdb, appConfig, logger)

	// instance service
	services := service.NewService(repo, appConfig, logger)

	// instance controller
	Ctl := handler.NewHandler(services, logger, jwtLib)

	mw := middleware.NewMiddleware(rdb, jwtLib)

	return &ServiceContext{Cacher: rdb, Cfg: appConfig, Ctl: *Ctl, Log: logger, Middleware: mw, JWT: jwtLib}, nil
}
