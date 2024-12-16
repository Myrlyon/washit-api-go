package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	orderRoutes "washit-api/app/order/routes"
	userRoutes "washit-api/app/user/routes"
	"washit-api/configs"
	dbs "washit-api/db"
	"washit-api/utils"
)

type Server struct {
	addr   string
	db     dbs.DatabaseInterface
	engine *gin.Engine
}

func NewServer(db dbs.DatabaseInterface) *Server {
	return &Server{
		addr:   configs.Envs.Port,
		db:     db,
		engine: gin.Default(),
	}
}

func (s *Server) Run() error {
	_ = s.engine.SetTrustedProxies(nil)

	s.engine.Use(gin.Recovery())
	s.engine.Use(gin.Logger())

	if err := s.MapRoutes(); err != nil {
		log.Fatalf("Mapping routes: %v", err)
	}

	// s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.engine.GET("/health", func(c *gin.Context) {
		utils.WriteJson(c, http.StatusOK, gin.H{"message": "OK"})
	})

	log.Println("HTTP server is listening on PORT: ", s.addr)
	if err := s.engine.Run(fmt.Sprintf(":%s", configs.Envs.Port)); err != nil {
		log.Fatalf("Running HTTP server: %v", err)
	}

	return nil
}

func (s Server) MapRoutes() error {
	v1 := s.engine.Group("/api/v1")
	userRoutes.Main(v1, s.db)
	orderRoutes.Main(v1, s.db)
	// productHttp.Routes(v1, s.db, s.validator, s.cache)
	// orderHttp.Routes(v1, s.db, s.validator)
	return nil
}

func (s Server) GetEngine() *gin.Engine {
	return s.engine
}
