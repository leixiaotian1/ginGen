package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/leixiaotian1/ginGen/internal/generator"
	"github.com/leixiaotian1/ginGen/internal/utils"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <feature> [project_path]",
	Short: "Add a feature (e.g., mysql) to an existing project",
	Args:  cobra.MinimumNArgs(1), // feature is required
	Run: func(cmd *cobra.Command, args []string) {
		feature := args[0]
		projectPath := "." // Default to current directory
		if len(args) > 1 {
			projectPath = args[1]
		}

		absProjectPath, err := filepath.Abs(projectPath)
		if err != nil {
			log.Fatalf("Error getting absolute path for %s: %v", projectPath, err)
		}

		// 1. Verify projectPath is a Go module
		goModPath := filepath.Join(absProjectPath, "go.mod")
		if _, err := os.Stat(goModPath); os.IsNotExist(err) {
			log.Fatalf("Error: %s is not a Go module (go.mod not found).\n", absProjectPath)
		}

		// Get module path from existing go.mod for template data
		currentModulePath, err := utils.GetModulePathFromGoMod(absProjectPath)
		if err != nil {
			log.Fatalf("Error reading module path from go.mod: %v", err)
		}
		templateData := generator.TemplateData{ModulePath: currentModulePath}

		fmt.Printf("Adding feature '%s' to project in '%s'...\n", feature, absProjectPath)

		switch feature {
		case "mysql", "gorm":
			addMySQLGormFeature(absProjectPath, templateData)
		case "postgres", "postgresql":
			addPostgresFeature(absProjectPath, templateData)
		case "redis":
			addRedisFeature(absProjectPath, templateData)
		case "kafka":
			addKafkaFeature(absProjectPath, templateData)
		case "jwt":
			addJWTFeature(absProjectPath, templateData)
		case "logger", "logging":
			addLoggerFeature(absProjectPath, templateData)
		case "swagger":
			addSwaggerFeature(absProjectPath, templateData)
		case "middleware":
			addMiddlewareFeature(absProjectPath, templateData)
		case "health":
			addHealthFeature(absProjectPath, templateData)
		case "prometheus":
			addPrometheusFeature(absProjectPath, templateData)
		case "hotreload", "hot-reload":
			addHotReloadFeature(absProjectPath, templateData)
		case "cron":
			addCronFeature(absProjectPath, templateData)
		case "upload":
			addUploadFeature(absProjectPath, templateData)
		default:
			log.Fatalf("Error: Unknown feature '%s'. Supported features: mysql, postgres, redis, kafka, jwt, logger, swagger, middleware, health, prometheus, hotreload, cron, upload.\n", feature)
		}

		fmt.Println("Running 'go mod tidy' to update dependencies...")
		if err := utils.RunCommand(absProjectPath, "go", "mod", "tidy"); err != nil {
			log.Printf("Warning: 'go mod tidy' failed: %v. Please run it manually.\n", err)
		}

		fmt.Printf("\nFeature '%s' added to project %s.\n", feature, absProjectPath)
		fmt.Println("Please review the generated files and update your main configuration and application logic as needed.")
	},
}

func addMySQLGormFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding MySQL (with GORM) support...")

	// 1. Go get GORM and MySQL driver
	pkgs := []string{
		"gorm.io/gorm",
		"gorm.io/driver/mysql",
		"github.com/go-sql-driver/mysql", // GORM's MySQL driver depends on this
	}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v. Please ensure it's added or try manually.\n", pkg, err)
			// For MVP, we'll continue, but in a real app, you might want to stop.
		}
	}

	// 2. Generate GORM client and DB config files from templates
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/mysql/config.go.tmpl", "internal/config/db_config.go"},
		{"templates/addfeature/mysql/client.go.tmpl", "internal/clients/gorm.go"},
	}

	// Ensure internal/clients directory exists
	if err := utils.CreateDirs(projectPath, "internal/clients"); err != nil {
		log.Fatalf("Failed to create internal/clients directory: %v", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data)
		if err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}

	// 3. Print message about adding to config.yaml and initializing in main.
	fmt.Println("\n--- Action Required ---")
	fmt.Println("1. Update 'configs/config.yaml' with your MySQL connection details.")
	fmt.Println("   Add a section like this (example):")

	// Print the content of config_entry.yaml.tmpl
	yamlSnippetPath := "templates/addfeature/mysql/config_entry.yaml.tmpl"
	yamlSnippet, err := generator.AllTemplatesFS.ReadFile(yamlSnippetPath)
	if err != nil {
		log.Printf("Warning: could not read YAML snippet template %s: %v\n", yamlSnippetPath, err)
		fmt.Print(`  db:
		mysql:
		  dsn: "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		  max_idle_conns: 10
		  max_open_conns: 100
		  conn_max_lifetime: 3600 # in seconds
	`)
		fmt.Println() // 单独添加一个换行
	} else {
		fmt.Println(string(yamlSnippet))
	}

	fmt.Println("2. Update 'internal/config/config.go':")
	fmt.Println("   - Add 'DBConfig `mapstructure:\"db\"`' to the `Config` struct.")
	fmt.Println("   - Ensure `LoadConfig()` loads this new section.")
	fmt.Println("3. Initialize the GORM client in 'cmd/server/main.go' or your application setup:")
	fmt.Println("   - Example: `db, err := clients.NewGORMClient(appConfig.DB.MySQL)`")
	fmt.Println("   - Pass the `*gorm.DB` instance to your services/handlers as needed.")
	fmt.Println("----------------------")
}

func addRedisFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding Redis support...")

	// 1. Go get Redis client
	pkgs := []string{
		"github.com/redis/go-redis/v9",
	}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v. Please ensure it's added or try manually.\n", pkg, err)
		}
	}

	// 2. Generate Redis client and config files from templates
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/redis/config.go.tmpl", "internal/config/redis_config.go"},
		{"templates/addfeature/redis/client.go.tmpl", "internal/clients/redis.go"},
	}

	// Ensure internal/clients directory exists
	if err := utils.CreateDirs(projectPath, "internal/clients"); err != nil {
		log.Fatalf("Failed to create internal/clients directory: %v", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data)
		if err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}

	// 3. Print instructions
	fmt.Println("\n--- Action Required ---")
	fmt.Println("1. Update 'configs/config.yaml' with your Redis connection details.")
	fmt.Println("   Add a section like this (example):")

	yamlSnippetPath := "templates/addfeature/redis/config_entry.yaml.tmpl"
	yamlSnippet, err := generator.AllTemplatesFS.ReadFile(yamlSnippetPath)
	if err != nil {
		log.Printf("Warning: could not read YAML snippet template %s: %v\n", yamlSnippetPath, err)
		fmt.Print(`  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
    pool_size: 10
    min_idle_conns: 5
`)
		fmt.Println()
	} else {
		fmt.Println(string(yamlSnippet))
	}

	fmt.Println("2. Update 'internal/config/config.go':")
	fmt.Println("   - Add 'Redis RedisConfig `mapstructure:\"redis\"`' to the `Config` struct.")
	fmt.Println("3. Initialize the Redis client in 'cmd/server/main.go' or your application setup:")
	fmt.Println("   - Example: `rdb, err := clients.NewRedisClient(appConfig.Redis)`")
	fmt.Println("   - Pass the `*redis.Client` instance to your services/handlers as needed.")
	fmt.Println("----------------------")
}

func addKafkaFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding Kafka support...")

	// 1. Go get Kafka client
	pkgs := []string{
		"github.com/IBM/sarama",
	}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v. Please ensure it's added or try manually.\n", pkg, err)
		}
	}

	// 2. Generate Kafka client and config files from templates
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/kafka/config.go.tmpl", "internal/config/kafka_config.go"},
		{"templates/addfeature/kafka/client.go.tmpl", "internal/clients/kafka.go"},
	}

	// Ensure internal/clients directory exists
	if err := utils.CreateDirs(projectPath, "internal/clients"); err != nil {
		log.Fatalf("Failed to create internal/clients directory: %v", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data)
		if err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}

	// 3. Print instructions
	fmt.Println("\n--- Action Required ---")
	fmt.Println("1. Update 'configs/config.yaml' with your Kafka connection details.")
	fmt.Println("   Add a section like this (example):")

	yamlSnippetPath := "templates/addfeature/kafka/config_entry.yaml.tmpl"
	yamlSnippet, err := generator.AllTemplatesFS.ReadFile(yamlSnippetPath)
	if err != nil {
		log.Printf("Warning: could not read YAML snippet template %s: %v\n", yamlSnippetPath, err)
		fmt.Print(`  kafka:
    brokers:
      - "localhost:9092"
    group_id: "my-consumer-group"
    client_id: "gin-app"
`)
		fmt.Println()
	} else {
		fmt.Println(string(yamlSnippet))
	}

	fmt.Println("2. Update 'internal/config/config.go':")
	fmt.Println("   - Add 'Kafka KafkaConfig `mapstructure:\"kafka\"`' to the `Config` struct.")
	fmt.Println("3. Initialize the Kafka producer/consumer in 'cmd/server/main.go' or your application setup:")
	fmt.Println("   - Example: `producer, err := clients.NewKafkaProducer(appConfig.Kafka)`")
	fmt.Println("   - Example: `consumer, err := clients.NewKafkaConsumer(appConfig.Kafka)`")
	fmt.Println("----------------------")
}

func addPostgresFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding PostgreSQL (with GORM) support...")
	pkgs := []string{"gorm.io/gorm", "gorm.io/driver/postgres"}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v\n", pkg, err)
		}
	}
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/postgres/config.go.tmpl", "internal/config/postgres_config.go"},
		{"templates/addfeature/postgres/client.go.tmpl", "internal/clients/postgres.go"},
	}
	if err := utils.CreateDirs(projectPath, "internal/clients"); err != nil {
		log.Fatalf("Failed to create internal/clients directory: %v", err)
	}
	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}
	fmt.Println("PostgreSQL support added successfully!")
}

func addJWTFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding JWT authentication support...")
	pkgs := []string{"github.com/golang-jwt/jwt/v5"}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v\n", pkg, err)
		}
	}
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
		log.Fatalf("Failed to create directories: %v", err)
	}
	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}
	fmt.Println("JWT authentication support added successfully!")
}

func addLoggerFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding structured logging support...")
	pkgs := []string{"go.uber.org/zap", "gopkg.in/natefinch/lumberjack.v2"}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v\n", pkg, err)
		}
	}
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/logger/config.go.tmpl", "internal/config/logger_config.go"},
		{"templates/addfeature/logger/logger.go.tmpl", "internal/logger/logger.go"},
		{"templates/addfeature/logger/middleware.go.tmpl", "internal/middleware/logger.go"},
	}
	if err := utils.CreateDirs(projectPath, "internal/logger", "internal/middleware"); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}
	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}
	fmt.Println("Structured logging support added successfully!")
}

func addSwaggerFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding Swagger documentation support...")
	pkgs := []string{"github.com/swaggo/swag/cmd/swag", "github.com/swaggo/gin-swagger", "github.com/swaggo/files"}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v\n", pkg, err)
		}
	}
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/swagger/swagger.go.tmpl", "internal/swagger/swagger.go"},
		{"templates/addfeature/swagger/main_swagger.go.tmpl", "docs/swagger.go"},
	}
	if err := utils.CreateDirs(projectPath, "internal/swagger", "docs"); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}
	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}
	fmt.Println("Swagger documentation support added successfully!")
	fmt.Println("Note: Run 'swag init -g cmd/server/main.go' to generate Swagger docs")
}

func addMiddlewareFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding common middleware...")
	pkgs := []string{"github.com/gin-contrib/cors", "github.com/gin-contrib/requestid"}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v\n", pkg, err)
		}
	}
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
		log.Fatalf("Failed to create internal/middleware directory: %v", err)
	}
	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}
	fmt.Println("Common middleware added successfully!")
}

func addHealthFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding health check endpoints...")
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/health/health.go.tmpl", "internal/handler/health.go"},
	}
	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}
	fmt.Println("Health check endpoints added successfully!")
}

func addPrometheusFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding Prometheus metrics support...")
	pkgs := []string{"github.com/prometheus/client_golang/prometheus", "github.com/prometheus/client_golang/prometheus/promhttp"}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v\n", pkg, err)
		}
	}
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/prometheus/metrics.go.tmpl", "internal/metrics/metrics.go"},
	}
	if err := utils.CreateDirs(projectPath, "internal/metrics"); err != nil {
		log.Fatalf("Failed to create internal/metrics directory: %v", err)
	}
	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}
	fmt.Println("Prometheus metrics support added successfully!")
}

func addHotReloadFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding config hot reload support...")
	pkgs := []string{"github.com/fsnotify/fsnotify"}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v\n", pkg, err)
		}
	}
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/hotreload/hotreload.go.tmpl", "internal/hotreload/hotreload.go"},
	}
	if err := utils.CreateDirs(projectPath, "internal/hotreload"); err != nil {
		log.Fatalf("Failed to create internal/hotreload directory: %v", err)
	}
	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}
	fmt.Println("Config hot reload support added successfully!")
}

func addCronFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding cron job scheduler support...")
	pkgs := []string{"github.com/robfig/cron/v3"}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v\n", pkg, err)
		}
	}
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/cron/cron.go.tmpl", "internal/cron/cron.go"},
	}
	if err := utils.CreateDirs(projectPath, "internal/cron"); err != nil {
		log.Fatalf("Failed to create internal/cron directory: %v", err)
	}
	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}
	fmt.Println("Cron job scheduler support added successfully!")
}

func addUploadFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding file upload support...")
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/upload/upload.go.tmpl", "internal/upload/upload.go"},
		{"templates/addfeature/upload/handler.go.tmpl", "internal/handler/upload.go"},
	}
	if err := utils.CreateDirs(projectPath, "internal/upload"); err != nil {
		log.Fatalf("Failed to create internal/upload directory: %v", err)
	}
	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		if err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data); err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}
	fmt.Println("File upload support added successfully!")
}

func init() {
	rootCmd.AddCommand(addCmd)
}
