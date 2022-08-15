package groot

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Host    string
	Port    int
	Cors    bool
	Log     *LogConfig
	OpenApi *OpenApiConfig
}

type Server struct {
	Engine *gin.Engine
	Cfg    *ServerConfig
}

func NewDefaultServer(cfg *ServerConfig) *Server {
	if cfg == nil {
		err := UnmarshalConfigByKey("server", cfg)
		if err != nil {
			panic(err)
		}
	}

	if cfg == nil {
		cfg = &ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
			Cors: false,
		}
	} else {
		if cfg.Host == "" {
			cfg.Host = "127.0.0.1"
		}
		if cfg.Port <= 0 {
			cfg.Port = 8080
		}
	}

	r := gin.Default()
	InitLogger(cfg.Log)
	r.Use(GinLogger(), GinRecovery(true))
	if cfg.Cors {
		r.Use(cors.Default())
	}
	InitOpenApi(cfg.OpenApi, r)
	// listenon := fmt.Sprintf("%s:%d", ip, port)
	// fmt.Printf("listen on %s", listenon)
	// return r
	return &Server{
		Engine: r,
	}
}

func (s *Server) Run() error {
	listenon := fmt.Sprintf("%s:%d", s.Cfg.Host, s.Cfg.Port)
	return s.Engine.Run(listenon)
}
