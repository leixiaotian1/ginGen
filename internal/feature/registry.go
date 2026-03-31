package feature

import (
	"fmt"
	"strings"
)

// Spec describes a feature for CLI/Web (single source of truth).
type Spec struct {
	ID          string
	Aliases     []string
	Name        string
	Description string
	Category    string
	apply       func(*Context) error
}

var byID map[string]*Spec
var aliasToCanonical map[string]string

func init() {
	list := []*Spec{
		{ID: "mysql", Aliases: []string{"gorm"}, Name: "MySQL / GORM", Description: "Add MySQL database support with GORM ORM", Category: "Database", apply: applyMySQL},
		{ID: "postgres", Aliases: []string{"postgresql"}, Name: "PostgreSQL / GORM", Description: "Add PostgreSQL database support with GORM ORM", Category: "Database", apply: applyPostgres},
		{ID: "redis", Name: "Redis", Description: "Add Redis cache and session support", Category: "Cache", apply: applyRedis},
		{ID: "kafka", Name: "Kafka", Description: "Add Apache Kafka message queue support", Category: "Message Queue", apply: applyKafka},
		{ID: "jwt", Name: "JWT Authentication", Description: "Add JWT-based authentication and authorization", Category: "Security", apply: applyJWT},
		{ID: "logger", Aliases: []string{"logging"}, Name: "Structured Logging", Description: "Add structured logging with zap and log rotation", Category: "Observability", apply: applyLogger},
		{ID: "swagger", Name: "Swagger/OpenAPI", Description: "Add Swagger documentation and API explorer", Category: "Documentation", apply: applySwagger},
		{ID: "middleware", Name: "Common Middleware", Description: "Add CORS, rate limiting, request ID, and recovery middleware", Category: "Middleware", apply: applyMiddleware},
		{ID: "health", Name: "Health Checks", Description: "Add health, readiness, and liveness check endpoints", Category: "Observability", apply: applyHealth},
		{ID: "prometheus", Name: "Prometheus Metrics", Description: "Add Prometheus metrics collection and /metrics endpoint", Category: "Observability", apply: applyPrometheus},
		{ID: "hotreload", Aliases: []string{"hot-reload"}, Name: "Config Hot Reload", Description: "Add configuration hot reload and graceful shutdown", Category: "Configuration", apply: applyHotReload},
		{ID: "cron", Name: "Task Scheduler", Description: "Add cron job scheduler for periodic tasks", Category: "Utilities", apply: applyCron},
		{ID: "upload", Name: "File Upload", Description: "Add file upload functionality with validation", Category: "Utilities", apply: applyUpload},
	}
	byID = make(map[string]*Spec, len(list))
	aliasToCanonical = make(map[string]string)
	for _, s := range list {
		byID[s.ID] = s
		aliasToCanonical[strings.ToLower(s.ID)] = s.ID
		for _, a := range s.Aliases {
			aliasToCanonical[strings.ToLower(a)] = s.ID
		}
	}
}

// Normalize returns canonical feature id or error.
func Normalize(raw string) (string, error) {
	key := strings.ToLower(strings.TrimSpace(raw))
	if key == "" {
		return "", fmt.Errorf("empty feature name")
	}
	id, ok := aliasToCanonical[key]
	if !ok {
		return "", fmt.Errorf("unknown feature %q", raw)
	}
	return id, nil
}

// AllSpecs returns registry order for API listing.
func AllSpecs() []*Spec {
	return []*Spec{
		byID["mysql"], byID["postgres"], byID["redis"], byID["kafka"],
		byID["jwt"], byID["logger"], byID["swagger"], byID["middleware"],
		byID["health"], byID["prometheus"], byID["hotreload"], byID["cron"], byID["upload"],
	}
}

// SupportedListForError is a comma-separated canonical ids string for help text.
func SupportedListForError() string {
	ids := []string{"mysql", "postgres", "redis", "kafka", "jwt", "logger", "swagger", "middleware", "health", "prometheus", "hotreload", "cron", "upload"}
	return strings.Join(ids, ", ")
}
