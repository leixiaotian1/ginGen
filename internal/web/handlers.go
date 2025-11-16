package web

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leixiaotian1/ginGen/internal/generator"
	"github.com/leixiaotian1/ginGen/internal/utils"
)

// ProjectRequest represents the request body for project generation
type ProjectRequest struct {
	ProjectName string   `json:"projectName" binding:"required"`
	ModulePath  string   `json:"modulePath" binding:"required"`
	Features    []string `json:"features"`
}

// AvailableFeatures returns the list of available features/modules
func GetAvailableFeatures(c *gin.Context) {
	features := []map[string]interface{}{
		{
			"id":          "mysql",
			"name":        "MySQL / GORM",
			"description": "Add MySQL database support with GORM ORM",
			"category":    "Database",
		},
		{
			"id":          "postgres",
			"name":        "PostgreSQL / GORM",
			"description": "Add PostgreSQL database support with GORM ORM",
			"category":    "Database",
		},
		{
			"id":          "redis",
			"name":        "Redis",
			"description": "Add Redis cache and session support",
			"category":    "Cache",
		},
		{
			"id":          "kafka",
			"name":        "Kafka",
			"description": "Add Apache Kafka message queue support",
			"category":    "Message Queue",
		},
		{
			"id":          "jwt",
			"name":        "JWT Authentication",
			"description": "Add JWT-based authentication and authorization",
			"category":    "Security",
		},
		{
			"id":          "logger",
			"name":        "Structured Logging",
			"description": "Add structured logging with zap and log rotation",
			"category":    "Observability",
		},
		{
			"id":          "swagger",
			"name":        "Swagger/OpenAPI",
			"description": "Add Swagger documentation and API explorer",
			"category":    "Documentation",
		},
		{
			"id":          "middleware",
			"name":        "Common Middleware",
			"description": "Add CORS, rate limiting, request ID, and recovery middleware",
			"category":    "Middleware",
		},
		{
			"id":          "health",
			"name":        "Health Checks",
			"description": "Add health, readiness, and liveness check endpoints",
			"category":    "Observability",
		},
		{
			"id":          "prometheus",
			"name":        "Prometheus Metrics",
			"description": "Add Prometheus metrics collection and /metrics endpoint",
			"category":    "Observability",
		},
		{
			"id":          "hotreload",
			"name":        "Config Hot Reload",
			"description": "Add configuration hot reload and graceful shutdown",
			"category":    "Configuration",
		},
		{
			"id":          "cron",
			"name":        "Task Scheduler",
			"description": "Add cron job scheduler for periodic tasks",
			"category":    "Utilities",
		},
		{
			"id":          "upload",
			"name":        "File Upload",
			"description": "Add file upload functionality with validation",
			"category":    "Utilities",
		},
	}

	c.JSON(200, gin.H{
		"features": features,
	})
}

// GenerateProject handles project generation request
func GenerateProject(c *gin.Context) {
	var req ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate project name
	if req.ProjectName == "" {
		c.JSON(400, gin.H{"error": "Project name is required"})
		return
	}

	// Validate module path
	if req.ModulePath == "" {
		req.ModulePath = req.ProjectName
	}

	// Create temporary directory for project generation
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("ginGen_%s_%d", req.ProjectName, time.Now().Unix()))
	defer os.RemoveAll(tempDir) // Clean up after

	projectPath := filepath.Join(tempDir, req.ProjectName)

	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		log.Printf("Error creating project directory: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create project directory"})
		return
	}

	// Generate base project structure
	templateData := generator.TemplateData{
		ProjectName: req.ProjectName,
		ModulePath:  req.ModulePath,
	}

	if err := generator.GenerateProjectStructure(projectPath, templateData); err != nil {
		log.Printf("Error generating project structure: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to generate project: %v", err)})
		return
	}

	// Add selected features
	featureMap := make(map[string]bool)
	for _, feature := range req.Features {
		featureMap[feature] = true
	}

	if featureMap["mysql"] {
		if err := addMySQLFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding MySQL feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add MySQL feature: %v", err)})
			return
		}
	}

	if featureMap["redis"] {
		if err := addRedisFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding Redis feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add Redis feature: %v", err)})
			return
		}
	}

	if featureMap["kafka"] {
		if err := addKafkaFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding Kafka feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add Kafka feature: %v", err)})
			return
		}
	}

	if featureMap["postgres"] {
		if err := addPostgresFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding PostgreSQL feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add PostgreSQL feature: %v", err)})
			return
		}
	}

	if featureMap["jwt"] {
		if err := addJWTFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding JWT feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add JWT feature: %v", err)})
			return
		}
	}

	if featureMap["logger"] {
		if err := addLoggerFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding Logger feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add Logger feature: %v", err)})
			return
		}
	}

	if featureMap["swagger"] {
		if err := addSwaggerFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding Swagger feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add Swagger feature: %v", err)})
			return
		}
	}

	if featureMap["middleware"] {
		if err := addMiddlewareFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding Middleware feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add Middleware feature: %v", err)})
			return
		}
	}

	if featureMap["health"] {
		if err := addHealthFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding Health feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add Health feature: %v", err)})
			return
		}
	}

	if featureMap["prometheus"] {
		if err := addPrometheusFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding Prometheus feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add Prometheus feature: %v", err)})
			return
		}
	}

	if featureMap["hotreload"] {
		if err := addHotReloadFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding HotReload feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add HotReload feature: %v", err)})
			return
		}
	}

	if featureMap["cron"] {
		if err := addCronFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding Cron feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add Cron feature: %v", err)})
			return
		}
	}

	if featureMap["upload"] {
		if err := addUploadFeature(projectPath, templateData); err != nil {
			log.Printf("Error adding Upload feature: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add Upload feature: %v", err)})
			return
		}
	}

	// Run go mod tidy (silently for web interface)
	if err := utils.RunCommandSilent(projectPath, "go", "mod", "tidy"); err != nil {
		log.Printf("Warning: go mod tidy failed: %v", err)
		// Continue anyway, the project is still generated
	}

	// Create ZIP file
	zipPath := filepath.Join(tempDir, fmt.Sprintf("%s.zip", req.ProjectName))
	if err := createZip(projectPath, zipPath); err != nil {
		log.Printf("Error creating ZIP file: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create project archive"})
		return
	}

	// Send ZIP file to client
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", req.ProjectName))
	c.Header("Content-Type", "application/zip")
	c.File(zipPath)
}

// addMySQLFeature adds MySQL/GORM support to the project
func addMySQLFeature(projectPath string, data generator.TemplateData) error {
	// Install dependencies (we'll note them in a file, actual installation happens via go mod tidy)
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/mysql/config.go.tmpl", "internal/config/db_config.go"},
		{"templates/addfeature/mysql/client.go.tmpl", "internal/clients/gorm.go"},
	}

	// Ensure internal/clients directory exists
	if err := utils.CreateDirs(projectPath, "internal/clients"); err != nil {
		return fmt.Errorf("failed to create internal/clients directory: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	// Update config.yaml with MySQL entry
	configYamlPath := filepath.Join(projectPath, "configs/config.yaml")
	if err := appendMySQLConfig(configYamlPath); err != nil {
		log.Printf("Warning: failed to append MySQL config to config.yaml: %v", err)
	}

	// Update config.go to include DBConfig
	configGoPath := filepath.Join(projectPath, "internal/config/config.go")
	if err := updateConfigGoForMySQL(configGoPath); err != nil {
		log.Printf("Warning: failed to update config.go for MySQL: %v", err)
	}

	return nil
}

// updateConfigGoForMySQL updates config.go to include DBConfig struct
func updateConfigGoForMySQL(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// Check if DBConfig already exists
	if strings.Contains(contentStr, "DBConfig") && strings.Contains(contentStr, "DB DBConfig") {
		return nil // Already updated
	}

	// Uncomment DBConfig in Config struct
	contentStr = strings.Replace(contentStr,
		"	// DB DBConfig `mapstructure:\"db\"` // Example for when DB is added",
		"	DB DBConfig `mapstructure:\"db\"`",
		1)

	// Uncomment DBConfig type definitions
	contentStr = strings.Replace(contentStr,
		"// DBConfig might be added by `ginGen add mysql`",
		"// DBConfig holds database configuration",
		1)

	contentStr = strings.Replace(contentStr,
		"// type DBConfig struct {",
		"type DBConfig struct {",
		1)

	contentStr = strings.Replace(contentStr,
		"//   MySQL MySQLConfig `mapstructure:\"mysql\"`",
		"	MySQL MySQLConfig `mapstructure:\"mysql\"`",
		1)

	contentStr = strings.Replace(contentStr,
		"// }",
		"}",
		1)

	contentStr = strings.Replace(contentStr,
		"// type MySQLConfig struct {",
		"type MySQLConfig struct {",
		1)

	contentStr = strings.Replace(contentStr,
		"// 	DSN             string `mapstructure:\"dsn\"`",
		"	DSN             string `mapstructure:\"dsn\"`",
		1)

	contentStr = strings.Replace(contentStr,
		"// 	MaxIdleConns    int    `mapstructure:\"max_idle_conns\"`",
		"	MaxIdleConns    int    `mapstructure:\"max_idle_conns\"`",
		1)

	contentStr = strings.Replace(contentStr,
		"// 	MaxOpenConns    int    `mapstructure:\"max_open_conns\"`",
		"	MaxOpenConns    int    `mapstructure:\"max_open_conns\"`",
		1)

	contentStr = strings.Replace(contentStr,
		"// 	ConnMaxLifetime int    `mapstructure:\"conn_max_lifetime\"` // in seconds",
		"	ConnMaxLifetime int    `mapstructure:\"conn_max_lifetime\"` // in seconds",
		1)

	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

// appendMySQLConfig appends MySQL configuration to config.yaml
func appendMySQLConfig(configPath string) error {
	yamlSnippetPath := "templates/addfeature/mysql/config_entry.yaml.tmpl"
	yamlSnippet, err := generator.AllTemplatesFS.ReadFile(yamlSnippetPath)
	if err != nil {
		return err
	}

	// Read existing config
	existingContent, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Append MySQL config
	newContent := string(existingContent) + "\n" + string(yamlSnippet)
	return os.WriteFile(configPath, []byte(newContent), 0644)
}

// addRedisFeature adds Redis support to the project
func addRedisFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/redis/config.go.tmpl", "internal/config/redis_config.go"},
		{"templates/addfeature/redis/client.go.tmpl", "internal/clients/redis.go"},
	}

	// Ensure internal/clients directory exists
	if err := utils.CreateDirs(projectPath, "internal/clients"); err != nil {
		return fmt.Errorf("failed to create internal/clients directory: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	// Update config.yaml with Redis entry
	configYamlPath := filepath.Join(projectPath, "configs/config.yaml")
	if err := appendRedisConfig(configYamlPath); err != nil {
		log.Printf("Warning: failed to append Redis config to config.yaml: %v", err)
	}

	// Update config.go to include RedisConfig
	configGoPath := filepath.Join(projectPath, "internal/config/config.go")
	if err := updateConfigGoForRedis(configGoPath); err != nil {
		log.Printf("Warning: failed to update config.go for Redis: %v", err)
	}

	return nil
}

// addKafkaFeature adds Kafka support to the project
func addKafkaFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/kafka/config.go.tmpl", "internal/config/kafka_config.go"},
		{"templates/addfeature/kafka/client.go.tmpl", "internal/clients/kafka.go"},
	}

	// Ensure internal/clients directory exists
	if err := utils.CreateDirs(projectPath, "internal/clients"); err != nil {
		return fmt.Errorf("failed to create internal/clients directory: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	// Update config.yaml with Kafka entry
	configYamlPath := filepath.Join(projectPath, "configs/config.yaml")
	if err := appendKafkaConfig(configYamlPath); err != nil {
		log.Printf("Warning: failed to append Kafka config to config.yaml: %v", err)
	}

	// Update config.go to include KafkaConfig
	configGoPath := filepath.Join(projectPath, "internal/config/config.go")
	if err := updateConfigGoForKafka(configGoPath); err != nil {
		log.Printf("Warning: failed to update config.go for Kafka: %v", err)
	}

	return nil
}

// appendRedisConfig appends Redis configuration to config.yaml
func appendRedisConfig(configPath string) error {
	yamlSnippetPath := "templates/addfeature/redis/config_entry.yaml.tmpl"
	yamlSnippet, err := generator.AllTemplatesFS.ReadFile(yamlSnippetPath)
	if err != nil {
		return err
	}

	existingContent, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	newContent := string(existingContent) + "\n" + string(yamlSnippet)
	return os.WriteFile(configPath, []byte(newContent), 0644)
}

// appendKafkaConfig appends Kafka configuration to config.yaml
func appendKafkaConfig(configPath string) error {
	yamlSnippetPath := "templates/addfeature/kafka/config_entry.yaml.tmpl"
	yamlSnippet, err := generator.AllTemplatesFS.ReadFile(yamlSnippetPath)
	if err != nil {
		return err
	}

	existingContent, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	newContent := string(existingContent) + "\n" + string(yamlSnippet)
	return os.WriteFile(configPath, []byte(newContent), 0644)
}

// updateConfigGoForRedis updates config.go to include RedisConfig
func updateConfigGoForRedis(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// Check if RedisConfig already exists in Config struct
	if strings.Contains(contentStr, "Redis RedisConfig") {
		return nil // Already updated
	}

	// Add Redis field to Config struct
	// Try to add after DB if exists, otherwise after Server
	if strings.Contains(contentStr, "DB DBConfig") {
		contentStr = strings.Replace(contentStr,
			"	DB DBConfig `mapstructure:\"db\"`",
			"	DB DBConfig `mapstructure:\"db\"`\n	Redis RedisConfig `mapstructure:\"redis\"`",
			1)
	} else {
		contentStr = strings.Replace(contentStr,
			"	Server ServerConfig `mapstructure:\"server\"`",
			"	Server ServerConfig `mapstructure:\"server\"`\n	Redis RedisConfig `mapstructure:\"redis\"`",
			1)
	}

	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

// updateConfigGoForKafka updates config.go to include KafkaConfig
func updateConfigGoForKafka(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// Check if KafkaConfig already exists in Config struct
	if strings.Contains(contentStr, "Kafka KafkaConfig") {
		return nil // Already updated
	}

	// Add Kafka field to Config struct - add after the last existing field
	// Priority: Redis > DB > Server
	if strings.Contains(contentStr, "Redis RedisConfig") {
		contentStr = strings.Replace(contentStr,
			"	Redis RedisConfig `mapstructure:\"redis\"`",
			"	Redis RedisConfig `mapstructure:\"redis\"`\n	Kafka KafkaConfig `mapstructure:\"kafka\"`",
			1)
	} else if strings.Contains(contentStr, "DB DBConfig") {
		contentStr = strings.Replace(contentStr,
			"	DB DBConfig `mapstructure:\"db\"`",
			"	DB DBConfig `mapstructure:\"db\"`\n	Kafka KafkaConfig `mapstructure:\"kafka\"`",
			1)
	} else {
		contentStr = strings.Replace(contentStr,
			"	Server ServerConfig `mapstructure:\"server\"`",
			"	Server ServerConfig `mapstructure:\"server\"`\n	Kafka KafkaConfig `mapstructure:\"kafka\"`",
			1)
	}

	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

// addPostgresFeature adds PostgreSQL support to the project
func addPostgresFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/postgres/config.go.tmpl", "internal/config/postgres_config.go"},
		{"templates/addfeature/postgres/client.go.tmpl", "internal/clients/postgres.go"},
	}

	if err := utils.CreateDirs(projectPath, "internal/clients"); err != nil {
		return fmt.Errorf("failed to create internal/clients directory: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	configYamlPath := filepath.Join(projectPath, "configs/config.yaml")
	if err := appendPostgresConfig(configYamlPath); err != nil {
		log.Printf("Warning: failed to append PostgreSQL config to config.yaml: %v", err)
	}

	configGoPath := filepath.Join(projectPath, "internal/config/config.go")
	if err := updateConfigGoForPostgres(configGoPath); err != nil {
		log.Printf("Warning: failed to update config.go for PostgreSQL: %v", err)
	}

	return nil
}

// addJWTFeature adds JWT authentication support
func addJWTFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/jwt/config.go.tmpl", "internal/config/jwt_config.go"},
		{"templates/addfeature/jwt/jwt.go.tmpl", "internal/utils/jwt/jwt.go"},
		{"templates/addfeature/jwt/middleware.go.tmpl", "internal/middleware/jwt_auth.go"},
		{"templates/addfeature/jwt/auth_handler.go.tmpl", "internal/handler/auth.go"},
	}

	if err := utils.CreateDirs(projectPath, "internal/utils/jwt", "internal/middleware"); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	configYamlPath := filepath.Join(projectPath, "configs/config.yaml")
	if err := appendJWTConfig(configYamlPath); err != nil {
		log.Printf("Warning: failed to append JWT config to config.yaml: %v", err)
	}

	configGoPath := filepath.Join(projectPath, "internal/config/config.go")
	if err := updateConfigGoForJWT(configGoPath); err != nil {
		log.Printf("Warning: failed to update config.go for JWT: %v", err)
	}

	return nil
}

// addLoggerFeature adds structured logging support
func addLoggerFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/logger/config.go.tmpl", "internal/config/logger_config.go"},
		{"templates/addfeature/logger/logger.go.tmpl", "internal/logger/logger.go"},
		{"templates/addfeature/logger/middleware.go.tmpl", "internal/middleware/logger.go"},
	}

	if err := utils.CreateDirs(projectPath, "internal/logger", "internal/middleware"); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	configYamlPath := filepath.Join(projectPath, "configs/config.yaml")
	if err := appendLoggerConfig(configYamlPath); err != nil {
		log.Printf("Warning: failed to append Logger config to config.yaml: %v", err)
	}

	configGoPath := filepath.Join(projectPath, "internal/config/config.go")
	if err := updateConfigGoForLogger(configGoPath); err != nil {
		log.Printf("Warning: failed to update config.go for Logger: %v", err)
	}

	return nil
}

// addSwaggerFeature adds Swagger documentation support
func addSwaggerFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/swagger/swagger.go.tmpl", "internal/swagger/swagger.go"},
		{"templates/addfeature/swagger/main_annotations.go.tmpl", "docs/swagger_annotations.txt"},
	}

	if err := utils.CreateDirs(projectPath, "internal/swagger", "docs"); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	return nil
}

// addMiddlewareFeature adds common middleware
func addMiddlewareFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/middleware/cors.go.tmpl", "internal/middleware/cors.go"},
		{"templates/addfeature/middleware/ratelimit.go.tmpl", "internal/middleware/ratelimit.go"},
		{"templates/addfeature/middleware/requestid.go.tmpl", "internal/middleware/requestid.go"},
		{"templates/addfeature/middleware/recovery.go.tmpl", "internal/middleware/recovery.go"},
	}

	if err := utils.CreateDirs(projectPath, "internal/middleware"); err != nil {
		return fmt.Errorf("failed to create internal/middleware directory: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	return nil
}

// addHealthFeature adds health check endpoints
func addHealthFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/health/health.go.tmpl", "internal/handler/health.go"},
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	return nil
}

// addPrometheusFeature adds Prometheus metrics
func addPrometheusFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/prometheus/metrics.go.tmpl", "internal/metrics/metrics.go"},
	}

	if err := utils.CreateDirs(projectPath, "internal/metrics"); err != nil {
		return fmt.Errorf("failed to create internal/metrics directory: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	return nil
}

// addHotReloadFeature adds config hot reload support
func addHotReloadFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/hotreload/hotreload.go.tmpl", "internal/hotreload/hotreload.go"},
	}

	if err := utils.CreateDirs(projectPath, "internal/hotreload"); err != nil {
		return fmt.Errorf("failed to create internal/hotreload directory: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	return nil
}

// addCronFeature adds cron job scheduler
func addCronFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/cron/cron.go.tmpl", "internal/cron/cron.go"},
	}

	if err := utils.CreateDirs(projectPath, "internal/cron"); err != nil {
		return fmt.Errorf("failed to create internal/cron directory: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	return nil
}

// addUploadFeature adds file upload functionality
func addUploadFeature(projectPath string, data generator.TemplateData) error {
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/upload/upload.go.tmpl", "internal/upload/upload.go"},
		{"templates/addfeature/upload/handler.go.tmpl", "internal/handler/upload.go"},
	}

	if err := utils.CreateDirs(projectPath, "internal/upload"); err != nil {
		return fmt.Errorf("failed to create internal/upload directory: %w", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			return fmt.Errorf("error generating file %s: %w", f.outputPath, err)
		}
	}

	return nil
}

// Helper functions for config updates
func appendPostgresConfig(configPath string) error {
	return appendConfigFromTemplate(configPath, "templates/addfeature/postgres/config_entry.yaml.tmpl")
}

func appendJWTConfig(configPath string) error {
	return appendConfigFromTemplate(configPath, "templates/addfeature/jwt/config_entry.yaml.tmpl")
}

func appendLoggerConfig(configPath string) error {
	return appendConfigFromTemplate(configPath, "templates/addfeature/logger/config_entry.yaml.tmpl")
}

func appendConfigFromTemplate(configPath, templatePath string) error {
	yamlSnippet, err := generator.AllTemplatesFS.ReadFile(templatePath)
	if err != nil {
		return err
	}

	existingContent, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	newContent := string(existingContent) + "\n" + string(yamlSnippet)
	return os.WriteFile(configPath, []byte(newContent), 0644)
}

func updateConfigGoForPostgres(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "PostgreSQL PostgreSQLConfig") {
		return nil
	}

	// Add PostgreSQL to DBConfig if exists, otherwise add as separate field
	if strings.Contains(contentStr, "DB DBConfig") {
		// PostgreSQL would be added to DBConfig struct
		return nil
	} else {
		// Add as separate field
		if strings.Contains(contentStr, "Redis RedisConfig") {
			contentStr = strings.Replace(contentStr,
				"	Redis RedisConfig `mapstructure:\"redis\"`",
				"	Redis RedisConfig `mapstructure:\"redis\"`\n	PostgreSQL PostgreSQLConfig `mapstructure:\"postgres\"`",
				1)
		} else if strings.Contains(contentStr, "DB DBConfig") {
			contentStr = strings.Replace(contentStr,
				"	DB DBConfig `mapstructure:\"db\"`",
				"	DB DBConfig `mapstructure:\"db\"`\n	PostgreSQL PostgreSQLConfig `mapstructure:\"postgres\"`",
				1)
		} else {
			contentStr = strings.Replace(contentStr,
				"	Server ServerConfig `mapstructure:\"server\"`",
				"	Server ServerConfig `mapstructure:\"server\"`\n	PostgreSQL PostgreSQLConfig `mapstructure:\"postgres\"`",
				1)
		}
	}

	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

func updateConfigGoForJWT(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "JWT JWTConfig") {
		return nil
	}

	// Add JWT field after the last existing field
	if strings.Contains(contentStr, "Kafka KafkaConfig") {
		contentStr = strings.Replace(contentStr,
			"	Kafka KafkaConfig `mapstructure:\"kafka\"`",
			"	Kafka KafkaConfig `mapstructure:\"kafka\"`\n	JWT JWTConfig `mapstructure:\"jwt\"`",
			1)
	} else if strings.Contains(contentStr, "Redis RedisConfig") {
		contentStr = strings.Replace(contentStr,
			"	Redis RedisConfig `mapstructure:\"redis\"`",
			"	Redis RedisConfig `mapstructure:\"redis\"`\n	JWT JWTConfig `mapstructure:\"jwt\"`",
			1)
	} else {
		contentStr = strings.Replace(contentStr,
			"	Server ServerConfig `mapstructure:\"server\"`",
			"	Server ServerConfig `mapstructure:\"server\"`\n	JWT JWTConfig `mapstructure:\"jwt\"`",
			1)
	}

	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

func updateConfigGoForLogger(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "Logger LoggerConfig") {
		return nil
	}

	// Add Logger field
	if strings.Contains(contentStr, "JWT JWTConfig") {
		contentStr = strings.Replace(contentStr,
			"	JWT JWTConfig `mapstructure:\"jwt\"`",
			"	JWT JWTConfig `mapstructure:\"jwt\"`\n	Logger LoggerConfig `mapstructure:\"logger\"`",
			1)
	} else {
		contentStr = strings.Replace(contentStr,
			"	Server ServerConfig `mapstructure:\"server\"`",
			"	Server ServerConfig `mapstructure:\"server\"`\n	Logger LoggerConfig `mapstructure:\"logger\"`",
			1)
	}

	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

// createZip creates a ZIP archive of the project directory
func createZip(sourceDir, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Calculate relative path for archive
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Use forward slashes in ZIP (Windows compatibility)
		relPath = strings.ReplaceAll(relPath, "\\", "/")

		// Create file in ZIP
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Copy file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		return err
	})
}
