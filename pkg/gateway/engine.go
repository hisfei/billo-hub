package gateway

import (
	"billohub/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewEngine initializes the Gin engine and configures the CORS middleware from the global configuration.
func NewEngine(debugMode bool) *gin.Engine {
	r := gin.New() // Use gin.New() to manually add middleware

	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configure CORS middleware from global config
	cfg := config.GetConfig()
	corsConfig := cors.DefaultConfig()

	// Allowed origins from config
	if len(cfg.CorsAllowedOrigins) > 0 {
		corsConfig.AllowOrigins = cfg.CorsAllowedOrigins
	} else {
		// If not specified in config, default to a secure empty list.
		corsConfig.AllowOrigins = []string{}
	}

	// Default allowed headers, can be customized further if needed
	corsConfig.AllowHeaders = []string{
		"Accept", "Content-Type", "DNT", "Funcid", "Referer", "Sec-Ch-Ua",
		"Sec-Ch-Ua-Mobile", "Sec-Ch-Ua-Platform", "Tenantid", "Tenantname",
		"Userid", "Username", "Authorization", "X-Requested-With", "X-Timestamp", "X-Nonce", "X-Sign", "X-UId",
	}

	// Default allowed methods
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	// Default exposed headers
	corsConfig.ExposeHeaders = []string{
		"Content-Length", "Funcid", "Tenantid", "Tenantname", "Userid", "Username", "X-Timestamp", "X-Nonce", "X-Sign", "X-UId",
	}

	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour // Set a reasonable default MaxAge

	r.Use(cors.New(corsConfig))

	return r
}
