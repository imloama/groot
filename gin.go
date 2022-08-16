package groot

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
	Cfg    *ServerConfig
}

func NewDefaultServer() *Server {
	LoadConfig()
	if serverCfg == nil {
		panic("error config file")
	}
	r := gin.Default()
	InitLogger(serverCfg.Log)
	r.Use(GinLogger(), GinRecovery(true))
	if serverCfg.Cors {
		r.Use(cors.Default())
	}
	// InitOpenApi(cfg.OpenApi, r)
	InitSwagger(serverCfg.Swagger, r)

	// listenon := fmt.Sprintf("%s:%d", ip, port)
	// fmt.Printf("listen on %s", listenon)
	// return r
	return &Server{
		Cfg:    serverCfg,
		Engine: r,
	}
}

func (s *Server) Run() error {
	listenon := fmt.Sprintf("%s:%d", s.Cfg.Host, s.Cfg.Port)
	return s.Engine.Run(listenon)
}