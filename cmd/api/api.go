package api

import (
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	orderRoutes "washit-api/app/order/routes"
	userRoutes "washit-api/app/user/routes"
	"washit-api/configs"
	dbs "washit-api/db"
	_ "washit-api/docs"
	"washit-api/redis"
	"washit-api/utils"
)

type Server struct {
	addr      string
	db        dbs.DatabaseInterface
	cache     redis.RedisInterface
	engine    *gin.Engine
	validator *validator.Validate
	app       *firebase.App
}

func NewServer(validator *validator.Validate, db dbs.DatabaseInterface, cache redis.RedisInterface, app *firebase.App) *Server {
	return &Server{
		addr:      configs.Envs.Port,
		db:        db,
		cache:     cache,
		engine:    gin.Default(),
		validator: validator,
		app:       app,
	}
}

func (s *Server) Run() error {
	_ = s.engine.SetTrustedProxies(nil)

	s.engine.Use(gin.Recovery())
	s.engine.Use(gin.Logger())

	if err := s.MapRoutes(); err != nil {
		log.Fatalf("Mapping routes: %v", err)
	}

	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.engine.GET("/ping", func(c *gin.Context) {
		utils.WriteJson(c, http.StatusOK, "pong")
	})

	log.Println("HTTP server is listening on PORT: ", s.addr)
	if err := s.engine.Run(fmt.Sprintf(":%s", configs.Envs.Port)); err != nil {
		log.Fatalf("Running HTTP server: %v", err)
	}

	return nil
}

func (s Server) MapRoutes() error {
	v1 := s.engine.Group("/api/v1")
	s.engine.Static("/public", "./public")
	userRoutes.Main(v1, s.db, s.cache, s.app, s.validator)
	orderRoutes.Main(v1, s.db, s.cache, s.validator)
	return nil
}

func (s Server) GetEngine() *gin.Engine {
	return s.engine
}
