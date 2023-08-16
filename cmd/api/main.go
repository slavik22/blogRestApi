package main

import (
	"context"
	"github.com/pkg/errors"
	"github.com/slavik22/blogRestApi"
	"github.com/slavik22/blogRestApi/controller"
	"github.com/slavik22/blogRestApi/lib/validator"
	"github.com/slavik22/blogRestApi/service"
	"github.com/slavik22/blogRestApi/store"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"

	echoLog "github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	// config
	cfg, err := blogRestApi.Get(".")

	// logger
	//l := logger.Get()

	//db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.MysqlUser, cfg.MysqlPassword, cfg.MysqlAddr, cfg.MysqlDB)))
	db, err := gorm.Open(mysql.Open(cfg.DBSource))

	// Init repository store (with mysql inside)
	store, err := store.New(ctx, db)

	if err != nil {
		return errors.Wrap(err, "store.New failed")
	}

	// Init service manager
	serviceManager, err := service.NewManager(ctx, store)
	if err != nil {
		return errors.Wrap(err, "manager.New failed")
	}

	// Init controllers
	userController := controller.NewUsers(ctx, serviceManager)

	// Initialize Echo instance
	e := echo.New()
	e.Validator = validator.NewValidator()
	//e.HTTPErrorHandler = libError.Error
	// Disable Echo JSON logger in debug mode
	if cfg.LogLevel == "debug" {
		if l, ok := e.Logger.(*echoLog.Logger); ok {
			l.SetHeader("${time_rfc3339} | ${level} | ${short_file}:${line}")
		}
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// API V1
	v1 := e.Group("/v1")

	// User routes
	userRoutes := v1.Group("/user")
	userRoutes.POST("/", userController.Create)
	userRoutes.GET("/:id", userController.Get)
	userRoutes.DELETE("/:id", userController.Delete)
	//userRoutes.PUT("/:id", userController.Update)

	// Start server
	s := &http.Server{
		Addr:         cfg.HTTPAddr,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))

	return nil
}
