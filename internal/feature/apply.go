package feature

import (
	"fmt"
	"path/filepath"

	"github.com/leixiaotian1/ginGen/internal/generator"
	"github.com/leixiaotian1/ginGen/internal/utils"
)

type tmplPair struct {
	templatePath string
	outputPath   string
}

func (c *Context) quiet() bool { return c.Opts.Quiet }

func (c *Context) runGoGet(pkgs []string) error {
	if c.Opts.SkipGoGet {
		return nil
	}
	for _, pkg := range pkgs {
		var err error
		if c.quiet() {
			err = utils.RunCommandSilent(c.ProjectPath, "go", "get", pkg)
		} else {
			err = utils.RunCommand(c.ProjectPath, "go", "get", pkg)
		}
		if err != nil {
			if c.Opts.Force {
				continue
			}
			return fmt.Errorf("go get %s: %w", pkg, err)
		}
	}
	return nil
}

func (c *Context) mkdir(dirs ...string) error {
	return utils.CreateDirsQuiet(c.ProjectPath, c.quiet(), dirs...)
}

func (c *Context) writeAll(pairs []tmplPair) error {
	fsys := c.TemplateFS()
	for _, f := range pairs {
		out := filepath.Join(c.ProjectPath, f.outputPath)
		if err := generator.CreateFileFromTemplateQuiet(fsys, f.templatePath, out, c.Data, c.quiet()); err != nil {
			return fmt.Errorf("generate %s: %w", f.outputPath, err)
		}
	}
	return nil
}

// Apply runs the feature pipeline for a canonical id (see registry).
func Apply(c *Context, id string) error {
	f := byID[id]
	if f == nil {
		return fmt.Errorf("unknown feature %q", id)
	}
	if !c.quiet() {
		fmt.Printf("Adding %s...\n", f.Name)
	}
	return f.apply(c)
}

func applyMySQL(c *Context) error {
	if err := c.runGoGet([]string{
		"gorm.io/gorm",
		"gorm.io/driver/mysql",
		"github.com/go-sql-driver/mysql",
	}); err != nil {
		return err
	}
	if err := c.mkdir("internal/clients"); err != nil {
		return err
	}
	if err := c.writeAll([]tmplPair{
		{"templates/addfeature/mysql/config.go.tmpl", "internal/config/db_config.go"},
		{"templates/addfeature/mysql/client.go.tmpl", "internal/clients/gorm.go"},
	}); err != nil {
		return err
	}
	return patchMySQL(c.TemplateFS(), c.ProjectPath)
}

func applyPostgres(c *Context) error {
	if err := c.runGoGet([]string{"gorm.io/gorm", "gorm.io/driver/postgres"}); err != nil {
		return err
	}
	if err := c.mkdir("internal/clients"); err != nil {
		return err
	}
	if err := c.writeAll([]tmplPair{
		{"templates/addfeature/postgres/config.go.tmpl", "internal/config/postgres_config.go"},
		{"templates/addfeature/postgres/client.go.tmpl", "internal/clients/postgres.go"},
	}); err != nil {
		return err
	}
	return patchPostgres(c.TemplateFS(), c.ProjectPath)
}

func applyRedis(c *Context) error {
	if err := c.runGoGet([]string{"github.com/redis/go-redis/v9"}); err != nil {
		return err
	}
	if err := c.mkdir("internal/clients"); err != nil {
		return err
	}
	if err := c.writeAll([]tmplPair{
		{"templates/addfeature/redis/config.go.tmpl", "internal/config/redis_config.go"},
		{"templates/addfeature/redis/client.go.tmpl", "internal/clients/redis.go"},
	}); err != nil {
		return err
	}
	return patchRedis(c.TemplateFS(), c.ProjectPath)
}

func applyKafka(c *Context) error {
	if err := c.runGoGet([]string{"github.com/IBM/sarama"}); err != nil {
		return err
	}
	if err := c.mkdir("internal/clients"); err != nil {
		return err
	}
	if err := c.writeAll([]tmplPair{
		{"templates/addfeature/kafka/config.go.tmpl", "internal/config/kafka_config.go"},
		{"templates/addfeature/kafka/client.go.tmpl", "internal/clients/kafka.go"},
	}); err != nil {
		return err
	}
	return patchKafka(c.TemplateFS(), c.ProjectPath)
}

func applyJWT(c *Context) error {
	if err := c.runGoGet([]string{"github.com/golang-jwt/jwt/v5"}); err != nil {
		return err
	}
	if err := c.mkdir("internal/utils/jwt", "internal/middleware"); err != nil {
		return err
	}
	if err := c.writeAll([]tmplPair{
		{"templates/addfeature/jwt/config.go.tmpl", "internal/config/jwt_config.go"},
		{"templates/addfeature/jwt/jwt.go.tmpl", "internal/utils/jwt/jwt.go"},
		{"templates/addfeature/jwt/middleware.go.tmpl", "internal/middleware/jwt_auth.go"},
		{"templates/addfeature/jwt/auth_handler.go.tmpl", "internal/handler/auth.go"},
	}); err != nil {
		return err
	}
	return patchJWT(c.TemplateFS(), c.ProjectPath)
}

func applyLogger(c *Context) error {
	if err := c.runGoGet([]string{"go.uber.org/zap", "gopkg.in/natefinch/lumberjack.v2"}); err != nil {
		return err
	}
	if err := c.mkdir("internal/logger", "internal/middleware"); err != nil {
		return err
	}
	if err := c.writeAll([]tmplPair{
		{"templates/addfeature/logger/config.go.tmpl", "internal/config/logger_config.go"},
		{"templates/addfeature/logger/logger.go.tmpl", "internal/logger/logger.go"},
		{"templates/addfeature/logger/middleware.go.tmpl", "internal/middleware/logger.go"},
	}); err != nil {
		return err
	}
	return patchLogger(c.TemplateFS(), c.ProjectPath)
}

func applySwagger(c *Context) error {
	if err := c.runGoGet([]string{
		"github.com/swaggo/swag/cmd/swag",
		"github.com/swaggo/gin-swagger",
		"github.com/swaggo/files",
	}); err != nil {
		return err
	}
	if err := c.mkdir("internal/swagger", "docs"); err != nil {
		return err
	}
	// Unified: CLI historical docs/swagger.go + Web annotations file
	if err := c.writeAll([]tmplPair{
		{"templates/addfeature/swagger/swagger.go.tmpl", "internal/swagger/swagger.go"},
		{"templates/addfeature/swagger/main_swagger.go.tmpl", "docs/swagger.go"},
		{"templates/addfeature/swagger/main_annotations.go.tmpl", "docs/swagger_annotations.txt"},
	}); err != nil {
		return err
	}
	return nil
}

func applyMiddleware(c *Context) error {
	if err := c.runGoGet([]string{"github.com/gin-contrib/cors", "github.com/gin-contrib/requestid"}); err != nil {
		return err
	}
	if err := c.mkdir("internal/middleware"); err != nil {
		return err
	}
	return c.writeAll([]tmplPair{
		{"templates/addfeature/middleware/cors.go.tmpl", "internal/middleware/cors.go"},
		{"templates/addfeature/middleware/ratelimit.go.tmpl", "internal/middleware/ratelimit.go"},
		{"templates/addfeature/middleware/requestid.go.tmpl", "internal/middleware/requestid.go"},
		{"templates/addfeature/middleware/recovery.go.tmpl", "internal/middleware/recovery.go"},
	})
}

func applyHealth(c *Context) error {
	return c.writeAll([]tmplPair{
		{"templates/addfeature/health/health.go.tmpl", "internal/handler/health.go"},
	})
}

func applyPrometheus(c *Context) error {
	if err := c.runGoGet([]string{
		"github.com/prometheus/client_golang/prometheus",
		"github.com/prometheus/client_golang/prometheus/promhttp",
	}); err != nil {
		return err
	}
	if err := c.mkdir("internal/metrics"); err != nil {
		return err
	}
	return c.writeAll([]tmplPair{
		{"templates/addfeature/prometheus/metrics.go.tmpl", "internal/metrics/metrics.go"},
	})
}

func applyHotReload(c *Context) error {
	if err := c.runGoGet([]string{"github.com/fsnotify/fsnotify"}); err != nil {
		return err
	}
	if err := c.mkdir("internal/hotreload"); err != nil {
		return err
	}
	return c.writeAll([]tmplPair{
		{"templates/addfeature/hotreload/hotreload.go.tmpl", "internal/hotreload/hotreload.go"},
	})
}

func applyCron(c *Context) error {
	if err := c.runGoGet([]string{"github.com/robfig/cron/v3"}); err != nil {
		return err
	}
	if err := c.mkdir("internal/cron"); err != nil {
		return err
	}
	return c.writeAll([]tmplPair{
		{"templates/addfeature/cron/cron.go.tmpl", "internal/cron/cron.go"},
	})
}

func applyUpload(c *Context) error {
	if err := c.mkdir("internal/upload"); err != nil {
		return err
	}
	return c.writeAll([]tmplPair{
		{"templates/addfeature/upload/upload.go.tmpl", "internal/upload/upload.go"},
		{"templates/addfeature/upload/handler.go.tmpl", "internal/handler/upload.go"},
	})
}
