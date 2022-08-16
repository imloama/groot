package groot

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/mvrilo/go-redoc"
// 	ginredoc "github.com/mvrilo/go-redoc/gin"
// )

// type OpenApiConfig struct {
// 	Enable      bool   `json:"enable"`
// 	Title       string `json:"title"`
// 	Description string `json:"descrption"`
// 	SpecFile    string `json:"spec_file"` //:    "./openapi.json", // "./openapi.yaml"
// 	SpecPath    string `json:"spec_path"` //:    "/openapi.json",  // "/openapi.yaml"
// 	DocsPath    string `json:"docs_path"` //:    "/docs",
// }

// func InitOpenApi(cfg *OpenApiConfig, r *gin.Engine) {
// 	if cfg == nil || !cfg.Enable {
// 		return
// 	}
// 	if cfg.SpecFile == "" {
// 		cfg.SpecFile = "./openapi.json"
// 	}
// 	if cfg.SpecPath == "" {
// 		cfg.SpecPath = "/openapi.json"
// 	}
// 	if cfg.DocsPath == "" {
// 		cfg.DocsPath = "/docs"
// 	}
// 	doc := redoc.Redoc{
// 		Title:       cfg.Title,
// 		Description: cfg.Description,
// 		SpecFile:    "./openapi.json", // "./openapi.yaml"
// 		SpecPath:    "/openapi.json",  // "/openapi.yaml"
// 		DocsPath:    "/docs",
// 	}
// 	r.Use(ginredoc.New(doc))
// }
