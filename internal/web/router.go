package web

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router for the web interface
func SetupRouter() *gin.Engine {
	// Set Gin to release mode for production (optional)
	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Enable CORS for web interface
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Determine static file path (works both in development and when installed)
	// Try multiple possible paths
	staticPaths := []string{
		"./static/web",                    // Development: from project root
		"static/web",                      // Alternative
		filepath.Join(".", "static", "web"), // Cross-platform
	}

	var staticPath string
	for _, path := range staticPaths {
		if _, err := os.Stat(path); err == nil {
			staticPath = path
			break
		}
	}

	// If no path found, try to find it relative to executable
	if staticPath == "" {
		exePath, err := os.Executable()
		if err == nil {
			exeDir := filepath.Dir(exePath)
			possiblePath := filepath.Join(exeDir, "static", "web")
			if _, err := os.Stat(possiblePath); err == nil {
				staticPath = possiblePath
			}
		}
	}

	// Serve static files
	if staticPath != "" {
		router.Static("/static", staticPath)
		router.StaticFile("/", filepath.Join(staticPath, "index.html"))
	} else {
		// Fallback: serve a simple message if static files not found
		router.GET("/", func(c *gin.Context) {
			c.String(500, "Static files not found. Please ensure static/web directory exists.")
		})
	}

	// API routes
	api := router.Group("/api")
	{
		api.GET("/features", GetAvailableFeatures)
		api.POST("/generate", GenerateProject)
	}

	return router
}

